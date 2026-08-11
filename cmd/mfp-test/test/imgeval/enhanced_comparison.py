import cv2
import numpy as np
from skimage.metrics import structural_similarity as ssim
from skimage.metrics import peak_signal_noise_ratio as psnr
from skimage.metrics import mean_squared_error
import matplotlib.pyplot as plt
from colormath.color_objects import sRGBColor, LabColor
from colormath.color_conversions import convert_color
from colormath.color_diff import delta_e_cie2000
import collections
import os
from scipy import ndimage
from scipy.stats import entropy

class ImageComparator:
    def __init__(self, original_path, processed_path, output_dir=None):
        """
        Initialize the comparator with original and processed images.
        
        Parameters:
        - original_path: Either a string path to an image file or a numpy array containing the image data
        - processed_path: Either a string path to an image file or a numpy array containing the image data
        - output_dir: Directory where output visualization images will be saved (default: current directory)
        """
        try:
            self.output_dir = output_dir if output_dir else "."
            
            if isinstance(original_path, str):
                self.original = cv2.imread(original_path)
                if self.original is None:
                    raise ValueError(f"Failed to load original image from {original_path}")
            else:
                self.original = original_path

            if isinstance(processed_path, str):
                self.processed = cv2.imread(processed_path)
                if self.processed is None:
                    raise ValueError(f"Failed to load processed image from {processed_path}")
            else:
                self.processed = processed_path

            if self.original is None or self.processed is None:
                raise ValueError("One or both images are None. Check input paths or arrays.")
            
            if self.original.shape != self.processed.shape:
                self.processed = cv2.resize(self.processed, 
                                            (self.original.shape[1], self.original.shape[0]))
            
            self.original_gray = cv2.cvtColor(self.original, cv2.COLOR_BGR2GRAY)
            self.processed_gray = cv2.cvtColor(self.processed, cv2.COLOR_BGR2GRAY)
            
            self.original_lab = cv2.cvtColor(self.original, cv2.COLOR_BGR2LAB)
            self.processed_lab = cv2.cvtColor(self.processed, cv2.COLOR_BGR2LAB)
        except Exception as e:
            print(f"Error initializing ImageComparator: {str(e)}")
            raise

    def basic_metrics(self):
        """Calculate basic image similarity metrics"""
        ssim_result = ssim(self.original_gray, self.processed_gray, full=True)
        ssim_score = ssim_result[0] if isinstance(ssim_result, tuple) else ssim_result
        
        mse_score = mean_squared_error(self.original_gray, self.processed_gray)
        
        if mse_score < 1e-10: 
            psnr_score = float('inf') 
        else:
            psnr_score = psnr(self.original_gray, self.processed_gray)
        
        return {
            "ssim": ssim_score,
            "mse": mse_score,
            "psnr": psnr_score
        }
    
    def feature_detection(self):
        """Detect and compare structural features between images"""
        
        orb = cv2.ORB_create()
        
        kp1, des1 = orb.detectAndCompute(self.original_gray, None)
        kp2, des2 = orb.detectAndCompute(self.processed_gray, None)
        
        # Check if descriptors are valid and have sufficient data
        if des1 is not None and des2 is not None and len(des1) > 0 and len(des2) > 0:
            bf = cv2.BFMatcher(cv2.NORM_HAMMING, crossCheck=True)
            
            try:
                matches = bf.match(des1, des2)
                
                matches = sorted(matches, key=lambda x: x.distance)
                
                if len(matches) > 0:
                    avg_distance = sum(m.distance for m in matches) / len(matches)
                    match_confidence = 1.0 - min(avg_distance / 100.0, 1.0)  # Normalize
                else:
                    match_confidence = 0.0
                    
                matched_img = cv2.drawMatches(self.original, kp1, self.processed, kp2, 
                                            matches[:30], None, flags=2)
                output_path = os.path.join(self.output_dir, 'feature_matches.png')
                cv2.imwrite(output_path, matched_img)
                
                return {
                    "keypoint_matches": len(matches),
                    "keypoints_original": len(kp1),
                    "keypoints_processed": len(kp2),
                    "match_confidence": match_confidence
                }
            except cv2.error as e:
                # Handle OpenCV errors (e.g., insufficient descriptors for matching)
                return {
                    "keypoint_matches": 0,
                    "keypoints_original": len(kp1) if kp1 else 0,
                    "keypoints_processed": len(kp2) if kp2 else 0,
                    "match_confidence": 0.0,
                    "error": f"Feature matching failed: {str(e)}"
                }
        else:
            return {
                "keypoint_matches": 0,
                "keypoints_original": len(kp1) if des1 is not None else 0,
                "keypoints_processed": len(kp2) if des2 is not None else 0,
                "match_confidence": 0.0
            }
    
    def edge_comparison(self):
        """Compare edges in the images"""
        
        edges1 = cv2.Canny(self.original_gray, 100, 200)
        edges2 = cv2.Canny(self.processed_gray, 100, 200)
        
        edge_result = ssim(edges1, edges2, full=True)
        edge_score = edge_result[0] if isinstance(edge_result, tuple) else edge_result
        
        cv2.imwrite(os.path.join(self.output_dir, 'edges_original.png'), edges1)
        cv2.imwrite(os.path.join(self.output_dir, 'edges_processed.png'), edges2)
        
        return {
            "edge_similarity": edge_score
        }
    
    def enhanced_edge_comparison(self):
        """Advanced edge detection and comparison between images"""
        
        original_canny = cv2.Canny(self.original_gray, 100, 200)
        processed_canny = cv2.Canny(self.processed_gray, 100, 200)
        
        original_canny_tight = cv2.Canny(self.original_gray, 150, 250)
        processed_canny_tight = cv2.Canny(self.processed_gray, 150, 250)
        
        original_canny_loose = cv2.Canny(self.original_gray, 50, 150)
        processed_canny_loose = cv2.Canny(self.processed_gray, 50, 150)
        
        edge_similarity = ssim(original_canny, processed_canny, full=True)[0]
        edge_similarity_tight = ssim(original_canny_tight, processed_canny_tight, full=True)[0]
        edge_similarity_loose = ssim(original_canny_loose, processed_canny_loose, full=True)[0]
        
        orig_edge_count = np.count_nonzero(original_canny)
        proc_edge_count = np.count_nonzero(processed_canny)
        
        if orig_edge_count > 0:
            edge_density_ratio = min(proc_edge_count / orig_edge_count, 2.0)  # Cap at 2.0 to prevent outliers
            edge_density_score = 1.0 - abs(1.0 - edge_density_ratio)
        else:
            edge_density_score = 0.0
        
        contours_orig, _ = cv2.findContours(original_canny, cv2.RETR_LIST, cv2.CHAIN_APPROX_SIMPLE)
        contours_proc, _ = cv2.findContours(processed_canny, cv2.RETR_LIST, cv2.CHAIN_APPROX_SIMPLE)
        
        orig_contour_lengths = [cv2.arcLength(cnt, False) for cnt in contours_orig]
        proc_contour_lengths = [cv2.arcLength(cnt, False) for cnt in contours_proc]
        
        orig_median_length = np.median(orig_contour_lengths) if len(orig_contour_lengths) > 0 else 0
        proc_median_length = np.median(proc_contour_lengths) if len(proc_contour_lengths) > 0 else 0
        
        if orig_median_length > 0:
            continuity_ratio = min(proc_median_length / orig_median_length, 2.0)
            edge_continuity_score = 1.0 - abs(1.0 - continuity_ratio)
        else:
            edge_continuity_score = 0.0
        
        cv2.imwrite(os.path.join(self.output_dir, 'edges_original.png'), original_canny)
        cv2.imwrite(os.path.join(self.output_dir, 'edges_processed.png'), processed_canny)
        cv2.imwrite(os.path.join(self.output_dir, 'edges_original_tight.png'), original_canny_tight)
        cv2.imwrite(os.path.join(self.output_dir, 'edges_processed_tight.png'), processed_canny_tight)
        
        vis_orig = cv2.cvtColor(self.original_gray, cv2.COLOR_GRAY2BGR)
        vis_proc = cv2.cvtColor(self.processed_gray, cv2.COLOR_GRAY2BGR)
        
        cv2.drawContours(vis_orig, contours_orig, -1, (0, 255, 0), 1)
        cv2.drawContours(vis_proc, contours_proc, -1, (0, 255, 0), 1)
        
        cv2.imwrite(os.path.join(self.output_dir, 'edges_original_contours.png'), vis_orig)
        cv2.imwrite(os.path.join(self.output_dir, 'edges_processed_contours.png'), vis_proc)
        
        edge_quality_score = (
            edge_similarity * 0.4 +
            (edge_similarity_tight + edge_similarity_loose) * 0.15 +
            edge_density_score * 0.15 +
            edge_continuity_score * 0.15
        )
        
        return {
            "edge_similarity": edge_similarity,
            "edge_similarity_tight": edge_similarity_tight,
            "edge_similarity_loose": edge_similarity_loose,
            "edge_density_score": edge_density_score,
            "edge_continuity_score": edge_continuity_score,
            "edge_quality_score": edge_quality_score
        }
    
    def color_analysis(self):
        """Analyze color differences between images"""
        hist_similarity = 0
        for i in range(3):
            hist1 = cv2.calcHist([self.original], [i], None, [256], [0, 256])
            hist2 = cv2.calcHist([self.processed], [i], None, [256], [0, 256])
            hist_similarity += cv2.compareHist(hist1, hist2, cv2.HISTCMP_CORREL)
        hist_similarity /= 3.0 
        
        # Calculate average Delta-E (color difference)
        height, width = self.original.shape[:2]
        sample_points = 1000
        delta_e_sum = 0
        count = 0
        
        for _ in range(sample_points):
            x = np.random.randint(0, width)
            y = np.random.randint(0, height)
            
            lab1 = self.original_lab[y, x]
            lab2 = self.processed_lab[y, x]
            
            color1 = LabColor(lab1[0], lab1[1], lab1[2])
            color2 = LabColor(lab2[0], lab2[1], lab2[2])
            
            try:
                delta_e = delta_e_cie2000(color1, color2)
                delta_e_sum += delta_e
                count += 1
            except Exception:
                pass
        
        avg_delta_e = delta_e_sum / count if count > 0 else float('inf')
        
        color_stats = {}
        for i, channel in enumerate(['B', 'G', 'R']):
            orig_mean = np.mean(self.original[:,:,i])
            proc_mean = np.mean(self.processed[:,:,i])
            orig_std = np.std(self.original[:,:,i])
            proc_std = np.std(self.processed[:,:,i])
            
            color_stats[f"{channel}_mean_diff"] = abs(orig_mean - proc_mean)
            color_stats[f"{channel}_std_diff"] = abs(orig_std - proc_std)
        
        return {
            "histogram_similarity": hist_similarity,
            "avg_delta_e": avg_delta_e,
            **color_stats
        }
    
    def generate_visual_diff(self):
        """Generate visual difference map between images"""
        
        diff = cv2.absdiff(self.original, self.processed)
        
        diff_enhanced = cv2.convertScaleAbs(diff, alpha=5.0)
        
        diff_heat = cv2.applyColorMap(diff_enhanced, cv2.COLORMAP_JET)
        
        cv2.imwrite(os.path.join(self.output_dir, 'visual_diff.png'), diff_heat)
        
        significant_diff = np.sum(diff > 20) / (diff.shape[0] * diff.shape[1] * diff.shape[2])
        return {
            "significant_diff_percent": significant_diff * 100
        }
    
    def page_integrity_comparison(self, full_document=False, pages_reference=None, pages_processed=None):
        """
        Verify page completeness and correct page ordering.
        
        Parameters:
        - full_document: If True, uses external page collections provided in pages_reference and pages_processed
        - pages_reference: List of reference page images (for multi-page documents)
        - pages_processed: List of processed page images to verify (for multi-page documents)
        
        Returns dictionary with:
        - page_count_match: Boolean indicating if page counts match
        - page_ordering_correct: Boolean indicating if pages are in correct order
        - correct_page_count: Number of expected pages
        - actual_page_count: Number of detected pages
        - page_completeness_score: Score between 0-1 indicating overall page integrity
        - missing_pages: List of missing page numbers (if any)
        - duplicated_pages: List of duplicated page numbers (if any)
        - reordered_pages: Dictionary mapping detected page to expected position
        """
        
        if full_document:
            if pages_reference is None or pages_processed is None:
                return {
                    "page_count_match": False,
                    "page_ordering_correct": False,
                    "page_completeness_score": 0,
                    "error": "Reference or processed page collections not provided"
                }
            
            ref_pages = pages_reference
            proc_pages = pages_processed
        else:
            return {
                "page_count_match": True, 
                "page_ordering_correct": True,
                "correct_page_count": 1,
                "actual_page_count": 1,
                "page_completeness_score": 1.0,
                "missing_pages": [],
                "duplicated_pages": []
            }
        
        correct_page_count = len(ref_pages)
        actual_page_count = len(proc_pages)
        page_count_match = correct_page_count == actual_page_count
        
        page_mapping = {}
        page_confidence = {}
        
        for i, proc_page in enumerate(proc_pages):
            best_match_idx = -1
            best_match_confidence = -1
            
            for j, ref_page in enumerate(ref_pages):
                if len(proc_page.shape) > 2:
                    proc_gray = cv2.cvtColor(proc_page, cv2.COLOR_BGR2GRAY)
                else:
                    proc_gray = proc_page
                    
                if len(ref_page.shape) > 2:
                    ref_gray = cv2.cvtColor(ref_page, cv2.COLOR_BGR2GRAY)
                else:
                    ref_gray = ref_page
                    
                orb = cv2.ORB_create()
                kp1, des1 = orb.detectAndCompute(ref_gray, None)
                kp2, des2 = orb.detectAndCompute(proc_gray, None)
                
                if des1 is not None and des2 is not None and len(des1) > 0 and len(des2) > 0 and len(kp1) > 0 and len(kp2) > 0:
                    try:
                        bf = cv2.BFMatcher(cv2.NORM_HAMMING, crossCheck=True)
                        matches = bf.match(des1, des2)
                        
                        if len(matches) > 0:
                            avg_distance = sum(m.distance for m in matches) / len(matches)
                            num_matches = len(matches)
                            confidence = num_matches / (1 + avg_distance/100)
                            
                            if confidence > best_match_confidence:
                                best_match_confidence = confidence
                                best_match_idx = j
                    except cv2.error:
                        # Skip matching if OpenCV encounters an error
                        pass
            
            if best_match_idx >= 0:
                page_mapping[i] = best_match_idx
                page_confidence[i] = best_match_confidence
        
        detected_pages = list(page_mapping.values())
        missing_pages = [i for i in range(correct_page_count) if i not in detected_pages]
        duplicated_pages = [page for page, count in collections.Counter(detected_pages).items() if count > 1]
        
        page_ordering_correct = True
        reordered_pages = {}
        
        for i in range(len(page_mapping) - 1):
            if i in page_mapping and i+1 in page_mapping:
                if page_mapping[i] >= page_mapping[i+1]:
                    page_ordering_correct = False
                    reordered_pages[i+1] = page_mapping[i+1]
        
        if correct_page_count == 0:
            page_completeness_score = 0
        else:
            missing_penalty = len(missing_pages) / correct_page_count if correct_page_count > 0 else 0
            duplicate_penalty = len(duplicated_pages) / correct_page_count if correct_page_count > 0 else 0
            reorder_penalty = len(reordered_pages) / correct_page_count if correct_page_count > 0 else 0
            
            page_completeness_score = 1.0 - (missing_penalty * 0.5 + duplicate_penalty * 0.3 + reorder_penalty * 0.2)
            page_completeness_score = max(0, min(1, page_completeness_score))  # Clip to [0,1]
        
        if len(proc_pages) > 0 and len(ref_pages) > 0:
            vis_height = 200
            vis_width = max(correct_page_count, actual_page_count) * 120
            visualization = np.ones((vis_height * 2 + 50, vis_width, 3), dtype=np.uint8) * 255
            
            for i in range(correct_page_count):
                left = i * 120 + 10
                cv2.rectangle(visualization, (left, 10), (left + 100, 160), (0, 0, 0), 2)
                cv2.putText(visualization, f"Ref #{i+1}", (left + 10, 40), 
                          cv2.FONT_HERSHEY_SIMPLEX, 0.5, (0, 0, 0), 1)
                
                if i < len(ref_pages):
                    page_img = ref_pages[i]
                    if len(page_img.shape) > 2:
                        thumb = cv2.resize(page_img, (80, 80))
                        visualization[60:140, left+10:left+90] = thumb
            
            for i in range(actual_page_count):
                left = i * 120 + 10
                cv2.rectangle(visualization, (left, vis_height + 10), (left + 100, vis_height + 160), 
                           (0, 0, 0), 2)
                cv2.putText(visualization, f"Proc #{i+1}", (left + 10, vis_height + 40), 
                          cv2.FONT_HERSHEY_SIMPLEX, 0.5, (0, 0, 0), 1)
                
                if i < len(proc_pages):
                    page_img = proc_pages[i]
                    if len(page_img.shape) > 2:
                        thumb = cv2.resize(page_img, (80, 80))
                        visualization[vis_height + 60:vis_height + 140, left+10:left+90] = thumb
                
                if i in page_mapping:
                    ref_idx = page_mapping[i]
                    ref_left = ref_idx * 120 + 60
                    proc_left = i * 120 + 60
                    color = (0, 255, 0) if ref_idx == i else (0, 0, 255)
                    cv2.line(visualization, (ref_left, 160), (proc_left, vis_height + 10), color, 2)
                    
                    if i in page_confidence:
                        conf_text = f"{page_confidence[i]:.2f}"
                        text_pos = ((ref_left + proc_left) // 2, (160 + vis_height + 10) // 2)
                        cv2.putText(visualization, conf_text, text_pos, 
                                  cv2.FONT_HERSHEY_SIMPLEX, 0.4, color, 1)
            
            cv2.imwrite(os.path.join(self.output_dir, 'page_integrity_map.png'), visualization)
        
        return {
            "page_count_match": page_count_match,
            "page_ordering_correct": page_ordering_correct,
            "correct_page_count": correct_page_count,
            "actual_page_count": actual_page_count,
            "page_completeness_score": page_completeness_score,
            "missing_pages": missing_pages,
            "duplicated_pages": duplicated_pages,
            "reordered_pages": reordered_pages,
            "page_mapping": page_mapping
        }
    
    def run_all_comparisons(self, full_document=False, pages_reference=None, pages_processed=None):
        """Run all comparison methods and return consolidated results"""
        results = {}
        results.update(self.basic_metrics())
        results.update(self.feature_detection())
        results.update(self.edge_comparison())
        results.update(self.enhanced_edge_comparison())
        results.update(self.color_analysis())
        results.update(self.generate_visual_diff())
        results.update(self.image_density_detection())
        results.update(self.rotation_detection())
        results.update(self.noise_analysis())
        results.update(self.texture_analysis())
        results.update(self.geometric_analysis())
        results.update(self.color_monochrome_detection())
        
        if full_document:
            results.update(self.page_integrity_comparison(True, pages_reference, pages_processed))
        
        psnr_factor = 1.0 if results["psnr"] == float('inf') else min(results["psnr"] / 50.0, 1.0)
        
        quality_score = (
            results["ssim"] * 0.15 + 
            (1.0 - min(results["mse"] / 10000.0, 1.0)) * 0.08 + 
            psnr_factor * 0.08 + 
            results["match_confidence"] * 0.12 + 
            results["edge_similarity"] * 0.04 +
            results["edge_quality_score"] * 0.12 +
            results["histogram_similarity"] * 0.08 +
            results["overall_content_density"] * 0.05 +
            (1.0 - min(abs(results["rotation_angle"]) / 45.0, 1.0)) * 0.05 +
            (1.0 - min(results["noise_level"] / 50.0, 1.0)) * 0.05 +
            results["texture_energy"] / 100.0 * 0.05 +
            results["avg_rectangularity"] * 0.05 +
            results["color_preservation_score"] * 0.04
        )
        
        if "page_completeness_score" in results:
            quality_score = quality_score * 0.85 + results["page_completeness_score"] * 0.15
        
        results["overall_quality"] = quality_score
        
        return results

    def image_density_detection(self):
        """
        Detect image density characteristics including:
        - Text density
        - Image content density
        - White space ratio
        - Content distribution
        """
        # Convert to grayscale and calculate overall pixel density (non-white content)
        gray = self.original_gray
        
        white_threshold = 240
        non_white_pixels = np.sum(gray < white_threshold)
        total_pixels = gray.shape[0] * gray.shape[1]
        content_density = non_white_pixels / total_pixels
        
        # Text density estimation using morphological operations
        kernel_text = cv2.getStructuringElement(cv2.MORPH_RECT, (2, 8))
        text_mask = cv2.morphologyEx(255 - gray, cv2.MORPH_CLOSE, kernel_text)
        text_density = np.sum(text_mask > 128) / total_pixels
        
        # Image content density (larger connected components)
        kernel_image = cv2.getStructuringElement(cv2.MORPH_RECT, (10, 10))
        image_mask = cv2.morphologyEx(255 - gray, cv2.MORPH_CLOSE, kernel_image)
        image_density = np.sum(image_mask > 128) / total_pixels

        white_space_ratio = 1.0 - content_density

        height, width = gray.shape
 
        grid_size = 16
        h_step = height // grid_size
        w_step = width // grid_size
        
        density_grid = np.zeros((grid_size, grid_size))
        for i in range(grid_size):
            for j in range(grid_size):
                y_start = i * h_step
                y_end = min((i + 1) * h_step, height)
                x_start = j * w_step
                x_end = min((j + 1) * w_step, width)
                
                region = gray[y_start:y_end, x_start:x_end]
                region_density = np.sum(region < white_threshold) / (region.shape[0] * region.shape[1])
                density_grid[i, j] = region_density
        
        # Calculate distribution uniformity
        density_variance = np.var(density_grid)
        density_entropy = entropy(density_grid.flatten() + 1e-10)
        
        # Detect content clusters
        dense_regions = density_grid > (content_density * 1.5)
        num_dense_clusters = self._count_connected_components(dense_regions)

        density_vis = cv2.applyColorMap((density_grid * 255).astype(np.uint8), cv2.COLORMAP_JET)
        density_vis = cv2.resize(density_vis, (width, height))
        cv2.imwrite(os.path.join(self.output_dir, 'density_heatmap.png'), density_vis)
        
        return {
            "overall_content_density": content_density,
            "text_density": text_density,
            "image_content_density": image_density,
            "white_space_ratio": white_space_ratio,
            "density_variance": density_variance,
            "density_entropy": density_entropy,
            "num_dense_clusters": num_dense_clusters,
            "density_grid": density_grid.tolist()
        }
    
    def rotation_detection(self):
        """
        Detect rotation angle of the image using multiple methods:
        - Hough line transform
        - Principal component analysis
        - Text line detection
        """
        gray = self.original_gray
        
        # Method 1: Hough Line Transform
        edges = cv2.Canny(gray, 50, 150, apertureSize=3)
        lines = cv2.HoughLines(edges, 1, np.pi/180, threshold=100)
        
        angles_hough = []
        if lines is not None:
            for line in lines:
                rho, theta = line[0]
                angle = theta * 180 / np.pi

                if angle > 90:
                    angle = angle - 180
                angles_hough.append(angle)
        
        # Method 2: Text line detection (for documents)
        angles_text = self._detect_text_rotation(gray)
        
        # Method 3: Principal Component Analysis
        angle_pca = self._detect_rotation_pca(gray)
        
        all_angles = angles_hough + angles_text
        if angle_pca is not None:
            all_angles.append(angle_pca)
        
        if all_angles:

            angles_array = np.array(all_angles)
            filtered_angles = angles_array[np.abs(angles_array) < 45]
            
            if len(filtered_angles) > 0:
                rotation_angle = np.median(filtered_angles)
                rotation_confidence = 1.0 - (np.std(filtered_angles) / 45.0)
                rotation_confidence = max(0, min(1, rotation_confidence))
            else:
                rotation_angle = 0
                rotation_confidence = 0
        else:
            rotation_angle = 0
            rotation_confidence = 0
        
        vis_img = cv2.cvtColor(gray, cv2.COLOR_GRAY2BGR)
        if lines is not None:
            for line in lines[:10]:
                rho, theta = line[0]
                a = np.cos(theta)
                b = np.sin(theta)
                x0 = a * rho
                y0 = b * rho
                x1 = int(x0 + 1000 * (-b))
                y1 = int(y0 + 1000 * (a))
                x2 = int(x0 - 1000 * (-b))
                y2 = int(y0 - 1000 * (a))
                cv2.line(vis_img, (x1, y1), (x2, y2), (0, 0, 255), 2)
        
        cv2.imwrite(os.path.join(self.output_dir, 'rotation_detection.png'), vis_img)
        
        return {
            "rotation_angle": rotation_angle,
            "rotation_confidence": rotation_confidence,
            "hough_angles": angles_hough,
            "text_angles": angles_text,
            "pca_angle": angle_pca,
            "is_rotated": abs(rotation_angle) > 1.0
        }
    
    def noise_analysis(self):
        """
        Analyze noise characteristics in the image:
        - Gaussian noise estimation
        - Salt and pepper noise
        - Noise distribution analysis
        """
        gray = self.original_gray.astype(np.float64)
        
        # Gaussian noise estimation using Laplacian
        laplacian = cv2.Laplacian(gray, cv2.CV_64F)
        noise_variance = laplacian.var()
        
        # Noise estimation using median filter
        gray_uint8 = self.original_gray.astype(np.uint8)
        median_filtered = cv2.medianBlur(gray_uint8, 5).astype(np.float64)
        noise_map = np.abs(gray - median_filtered)
        noise_level = np.mean(noise_map)
        
        # Salt and pepper noise detection
        kernel = np.ones((3, 3), np.uint8)
        salt_pepper_mask = np.zeros_like(gray)
        
        # Salt noise (isolated bright pixels)
        bright_thresh = np.percentile(gray, 95)
        bright_pixels = gray > bright_thresh
        eroded_bright = cv2.erode(bright_pixels.astype(np.uint8), kernel)
        salt_noise = bright_pixels & (eroded_bright == 0)
        
        # Pepper noise (isolated dark pixels)
        dark_thresh = np.percentile(gray, 5)
        dark_pixels = gray < dark_thresh
        eroded_dark = cv2.erode(dark_pixels.astype(np.uint8), kernel)
        pepper_noise = dark_pixels & (eroded_dark == 0)
        
        salt_count = np.sum(salt_noise)
        pepper_count = np.sum(pepper_noise)
        total_pixels = gray.shape[0] * gray.shape[1]
        
        # Noise distribution analysis
        noise_histogram, bins = np.histogram(noise_map, bins=50, range=(0, 255))
        noise_entropy = entropy(noise_histogram + 1e-10)
        
        # Signal-to-noise ratio estimation
        signal_power = np.var(gray)
        if noise_variance > 0:
            snr = 10 * np.log10(signal_power / noise_variance)
        else:
            snr = float('inf')
        
        noise_vis = cv2.applyColorMap((noise_map * 10).astype(np.uint8), cv2.COLORMAP_HOT)
        cv2.imwrite(os.path.join(self.output_dir, 'noise_analysis.png'), noise_vis)
        
        return {
            "noise_variance": float(noise_variance),
            "noise_level": float(noise_level),
            "salt_noise_count": int(salt_count),
            "pepper_noise_count": int(pepper_count),
            "salt_pepper_ratio": float((salt_count + pepper_count) / total_pixels),
            "noise_entropy": float(noise_entropy),
            "snr_db": float(snr),
            "noise_quality": "low" if noise_level > 10 else "medium" if noise_level > 5 else "high"
        }
    
    def texture_analysis(self):
        """
        Analyze texture characteristics using various methods:
        - Local Binary Patterns
        - Contrast measures
        - Texture uniformity
        """
        gray = self.original_gray
        
        # Local Binary Pattern
        lbp = self._calculate_lbp(gray)
        lbp_histogram, _ = np.histogram(lbp, bins=256, range=(0, 256))
        lbp_uniformity = np.sum(lbp_histogram ** 2) / (np.sum(lbp_histogram) ** 2)
        
        # Contrast measures
        kernel_size = 9
        kernel = np.ones((kernel_size, kernel_size)) / (kernel_size ** 2)
        local_mean = cv2.filter2D(gray.astype(np.float32), -1, kernel)
        local_variance = cv2.filter2D((gray.astype(np.float32) - local_mean) ** 2, -1, kernel)
        local_std = np.sqrt(local_variance)
        
        contrast_measure = np.mean(local_std)
        
        # Texture energy and homogeneity using co-occurrence matrix
        grad_x = cv2.Sobel(gray, cv2.CV_64F, 1, 0, ksize=3)
        grad_y = cv2.Sobel(gray, cv2.CV_64F, 0, 1, ksize=3)
        
        gradient_magnitude = np.sqrt(grad_x**2 + grad_y**2)
        texture_energy = np.mean(gradient_magnitude)
        
        # Texture directionality
        gradient_direction = np.arctan2(grad_y, grad_x)
        direction_histogram, _ = np.histogram(gradient_direction, bins=36, range=(-np.pi, np.pi))
        direction_uniformity = 1.0 - (entropy(direction_histogram + 1e-10) / np.log(36))
        
        # Roughness measure
        roughness = np.std(gray)
        
        texture_vis = cv2.applyColorMap(lbp, cv2.COLORMAP_JET)
        cv2.imwrite(os.path.join(self.output_dir, 'texture_lbp.png'), texture_vis)
        
        contrast_vis = cv2.applyColorMap((local_std * 255 / np.max(local_std)).astype(np.uint8), cv2.COLORMAP_HOT)
        cv2.imwrite(os.path.join(self.output_dir, 'texture_contrast.png'), contrast_vis)
        
        return {
            "lbp_uniformity": float(lbp_uniformity),
            "contrast_measure": float(contrast_measure),
            "texture_energy": float(texture_energy),
            "direction_uniformity": float(direction_uniformity),
            "roughness": float(roughness),
            "texture_complexity": "high" if texture_energy > 50 else "medium" if texture_energy > 20 else "low"
        }
    
    def geometric_analysis(self):
        """
        Analyze geometric properties and distortions:
        - Aspect ratio
        - Geometric distortions
        - Symmetry analysis
        """
        gray = self.original_gray
        height, width = gray.shape
        
        aspect_ratio = width / height
        
        # Detect rectangular structures for distortion analysis
        edges = cv2.Canny(gray, 50, 150)
        contours, _ = cv2.findContours(edges, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
        
        rectangularity_scores = []
        for contour in contours:
            if cv2.contourArea(contour) > 1000:  # Filter small contours
                epsilon = 0.02 * cv2.arcLength(contour, True)
                approx = cv2.approxPolyDP(contour, epsilon, True)
                
                if len(approx) == 4:
                    rect = cv2.minAreaRect(contour)
                    box = cv2.boxPoints(rect)
                    box = np.array(box, dtype=np.int32)
                    
                    rect_area = rect[1][0] * rect[1][1]
                    contour_area = cv2.contourArea(contour)
                    
                    if rect_area > 0:
                        rectangularity = contour_area / rect_area
                        rectangularity_scores.append(rectangularity)
        
        avg_rectangularity = np.mean(rectangularity_scores) if rectangularity_scores else 0
        
        # Symmetry analysis
        # Horizontal symmetry
        top_half = gray[:height//2, :]
        bottom_half = gray[height//2:, :]
        bottom_half_flipped = np.flipud(bottom_half[:top_half.shape[0], :])
        
        if top_half.shape == bottom_half_flipped.shape:
            h_sym_result = ssim(top_half, bottom_half_flipped)
            horizontal_symmetry = h_sym_result[0] if isinstance(h_sym_result, tuple) else h_sym_result
        else:
            horizontal_symmetry = 0.0
        
        # Vertical symmetry
        left_half = gray[:, :width//2]
        right_half = gray[:, width//2:]
        right_half_flipped = np.fliplr(right_half[:, :left_half.shape[1]])
        
        if left_half.shape == right_half_flipped.shape:
            v_sym_result = ssim(left_half, right_half_flipped)
            vertical_symmetry = v_sym_result[0] if isinstance(v_sym_result, tuple) else v_sym_result
        else:
            vertical_symmetry = 0.0
        
        # Line straightness analysis
        lines = cv2.HoughLinesP(edges, 1, np.pi/180, threshold=50, minLineLength=100, maxLineGap=10)
        
        line_straightness_scores = []
        if lines is not None:
            for line in lines:
                x1, y1, x2, y2 = line[0]
                line_length = np.sqrt((x2-x1)**2 + (y2-y1)**2)
                if line_length > 0:
                    # Perfect line would have straightness score of 1
                    straightness = 1.0  # Simplified - could be enhanced with curve fitting
                    line_straightness_scores.append(straightness)
        
        avg_line_straightness = np.mean(line_straightness_scores) if line_straightness_scores else 0
        
        return {
            "aspect_ratio": float(aspect_ratio),
            "avg_rectangularity": float(avg_rectangularity),
            "horizontal_symmetry": float(horizontal_symmetry),
            "vertical_symmetry": float(vertical_symmetry),
            "avg_line_straightness": float(avg_line_straightness),
            "num_rectangular_objects": len(rectangularity_scores),
            "geometric_quality": "high" if avg_rectangularity > 0.8 else "medium" if avg_rectangularity > 0.6 else "low"
        }
    
    def color_monochrome_detection(self):
        """
        Detect whether images are color or monochrome (grayscale).
        This method can distinguish RGB-encoded grayscale from RGB-encoded color images.
        
        Returns detailed analysis of color content for both original and processed images.
        """
        def analyze_image_color_content(image, name):
            """Analyze color content of a single image"""
            height, width = image.shape[:2]
            
            # Method 1: Channel variance analysis
            # For grayscale images encoded as RGB, all channels should be identical
            b_channel = image[:, :, 0].astype(np.float32)
            g_channel = image[:, :, 1].astype(np.float32)
            r_channel = image[:, :, 2].astype(np.float32)
            
            # Calculate channel differences
            bg_diff = np.abs(b_channel - g_channel)
            br_diff = np.abs(b_channel - r_channel)
            gr_diff = np.abs(g_channel - r_channel)
            
            # Statistics of channel differences
            max_bg_diff = np.max(bg_diff)
            max_br_diff = np.max(br_diff)
            max_gr_diff = np.max(gr_diff)
            
            mean_bg_diff = np.mean(bg_diff)
            mean_br_diff = np.mean(br_diff)
            mean_gr_diff = np.mean(gr_diff)
            
            # Overall channel difference metrics
            max_channel_diff = max(max_bg_diff, max_br_diff, max_gr_diff)
            mean_channel_diff = (mean_bg_diff + mean_br_diff + mean_gr_diff) / 3.0
            
            # Method 2: Color saturation analysis
            hsv = cv2.cvtColor(image, cv2.COLOR_BGR2HSV)
            saturation = hsv[:, :, 1].astype(np.float32)
            
            # Saturation statistics
            mean_saturation = np.mean(saturation)
            max_saturation = np.max(saturation)
            std_saturation = np.std(saturation)
            
            # Count pixels with significant saturation (indicating color)
            color_threshold = 15  # Threshold for considering a pixel as colored
            color_pixels = np.sum(saturation > color_threshold)
            color_pixel_ratio = color_pixels / (height * width)
            
            # Method 3: RGB channel correlation analysis
            # For grayscale, all channels should be highly correlated
            bg_correlation = np.corrcoef(b_channel.flatten(), g_channel.flatten())[0, 1]
            br_correlation = np.corrcoef(b_channel.flatten(), r_channel.flatten())[0, 1]
            gr_correlation = np.corrcoef(g_channel.flatten(), r_channel.flatten())[0, 1]
            
            min_correlation = min(bg_correlation, br_correlation, gr_correlation)
            mean_correlation = (bg_correlation + br_correlation + gr_correlation) / 3.0
            
            # Method 4: Color variance in Lab color space
            lab = cv2.cvtColor(image, cv2.COLOR_BGR2LAB)
            a_channel = lab[:, :, 1].astype(np.float32)  # Green-Red axis
            b_channel_lab = lab[:, :, 2].astype(np.float32)  # Blue-Yellow axis
            
            a_variance = np.var(a_channel)
            b_variance = np.var(b_channel_lab)
            color_variance = a_variance + b_variance
            
            # Method 5: Distinct color counting
            # Reshape image for unique color counting
            reshaped = image.reshape(-1, 3)
            unique_colors = len(np.unique(reshaped.view(np.dtype((np.void, reshaped.dtype.itemsize * reshaped.shape[1])))))
            
            # For true grayscale, we expect R=G=B for all pixels
            identical_channels = np.sum((image[:,:,0] == image[:,:,1]) & 
                                       (image[:,:,1] == image[:,:,2]) & 
                                       (image[:,:,0] == image[:,:,2]))
            identical_channel_ratio = identical_channels / (height * width)
            
            # Decision thresholds
            # An image is considered grayscale if:
            # 1. Channel differences are minimal
            # 2. Saturation is low
            # 3. Channel correlations are high
            # 4. Color variance in Lab space is low
            # 5. High ratio of identical RGB values
            
            is_grayscale_channel_diff = max_channel_diff < 5 and mean_channel_diff < 1
            is_grayscale_saturation = mean_saturation < 10 and color_pixel_ratio < 0.05
            is_grayscale_correlation = min_correlation > 0.99
            is_grayscale_lab_variance = color_variance < 100
            is_grayscale_identical = identical_channel_ratio > 0.95
            
            # Confidence scoring
            confidence_scores = []
            
            if max_channel_diff < 1:
                confidence_scores.append(1.0)
            elif max_channel_diff < 5:
                confidence_scores.append(0.8)
            else:
                confidence_scores.append(0.0)
                
            if mean_saturation < 5:
                confidence_scores.append(1.0)
            elif mean_saturation < 15:
                confidence_scores.append(0.6)
            else:
                confidence_scores.append(0.0)
                
            if min_correlation > 0.995:
                confidence_scores.append(1.0)
            elif min_correlation > 0.98:
                confidence_scores.append(0.7)
            else:
                confidence_scores.append(0.0)
                
            confidence = np.mean(confidence_scores)
            
            # Final decision
            votes = [is_grayscale_channel_diff, is_grayscale_saturation, 
                    is_grayscale_correlation, is_grayscale_lab_variance, 
                    is_grayscale_identical]
            
            is_grayscale = sum(votes) >= 3
            
            return {
                f"{name}_is_grayscale": is_grayscale,
                f"{name}_confidence": confidence,
                f"{name}_max_channel_diff": max_channel_diff,
                f"{name}_mean_saturation": mean_saturation,
                f"{name}_color_pixel_ratio": color_pixel_ratio,
                f"{name}_min_correlation": min_correlation,
                f"{name}_mean_correlation": mean_correlation,
                f"{name}_color_variance_lab": color_variance,
                f"{name}_unique_colors": unique_colors,
                f"{name}_identical_channel_ratio": identical_channel_ratio,
                f"{name}_votes": {
                    "channel_diff": is_grayscale_channel_diff,
                    "saturation": is_grayscale_saturation,
                    "correlation": is_grayscale_correlation,
                    "lab_variance": is_grayscale_lab_variance,
                    "identical_channels": is_grayscale_identical
                }
            }
        
        # Analyze both images
        original_results = analyze_image_color_content(self.original, "original")
        processed_results = analyze_image_color_content(self.processed, "processed")
        
        color_consistency = original_results["original_is_grayscale"] == processed_results["processed_is_grayscale"]
        
        # overall assessment
        if original_results["original_is_grayscale"] and processed_results["processed_is_grayscale"]:
            color_type_assessment = "Both images are grayscale"
        elif not original_results["original_is_grayscale"] and not processed_results["processed_is_grayscale"]:
            color_type_assessment = "Both images are color"
        elif original_results["original_is_grayscale"] and not processed_results["processed_is_grayscale"]:
            color_type_assessment = "Original is grayscale, processed is color (colorization occurred)"
        else:
            color_type_assessment = "Original is color, processed is grayscale (desaturization occurred)"
        
        # Calculate color preservation score
        if color_consistency:
            color_preservation_score = 1.0
        else:
            # Penalize color type changes
            color_preservation_score = 0.5
        
        # Combine all results
        results = {
            **original_results,
            **processed_results,
            "color_consistency": color_consistency,
            "color_type_assessment": color_type_assessment,
            "color_preservation_score": color_preservation_score
        }
        
        return results

    def _count_connected_components(self, binary_mask):
        """Count connected components in binary mask"""
        num_labels, _ = cv2.connectedComponents(binary_mask.astype(np.uint8))
        return num_labels - 1  # Subtract background
    
    def _detect_text_rotation(self, gray):
        """Detect rotation based on text line orientation"""
        # Detect horizontal text lines
        kernel = cv2.getStructuringElement(cv2.MORPH_RECT, (20, 1))
        horizontal_lines = cv2.morphologyEx(255 - gray, cv2.MORPH_OPEN, kernel)
        
        # Find contours of text lines
        contours, _ = cv2.findContours(horizontal_lines, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
        
        angles = []
        for contour in contours:
            if cv2.contourArea(contour) > 100:  # small contours
                rect = cv2.minAreaRect(contour)
                angle = rect[2]
                if angle < -45:
                    angle = 90 + angle
                angles.append(angle)
        
        return angles
    
    def _detect_rotation_pca(self, gray):
        """Detect rotation using Principal Component Analysis"""
        # edge points
        edges = cv2.Canny(gray, 50, 150)
        points = np.column_stack(np.where(edges > 0))
        
        if len(points) < 10:
            return None
        
        # Calculate PCA
        mean = np.mean(points, axis=0)
        centered_points = points - mean
        
        # covariance matrix
        cov_matrix = np.cov(centered_points.T)
        eigenvalues, eigenvectors = np.linalg.eig(cov_matrix)
        
        # Principal component (maximum variance)
        principal_component = eigenvectors[:, np.argmax(eigenvalues)]
        angle = np.arctan2(principal_component[1], principal_component[0]) * 180 / np.pi
        
        # Convert to rotation from horizontal
        if angle > 45:
            angle = angle - 90
        elif angle < -45:
            angle = angle + 90
        
        return angle
    
    def _calculate_lbp(self, gray, radius=1, neighbors=8):
        """Calculate Local Binary Pattern (optimized version)"""
        # Use a much smaller sample for speed - downsample the image
        height, width = gray.shape
        if height > 256 or width > 256:
            # Downsample for faster processing
            scale_factor = min(256 / height, 256 / width)
            new_height = int(height * scale_factor)
            new_width = int(width * scale_factor)
            gray = cv2.resize(gray, (new_width, new_height))
        
        lbp = np.zeros_like(gray)
        
        # Simplified LBP calculation using vectorized operations
        h, w = gray.shape
        for i in range(radius, h - radius, 4):  # Skip every 4th pixel for speed
            for j in range(radius, w - radius, 4):  # Skip every 4th pixel for speed
                center = gray[i, j]
                
                # Sample only 4 neighbors instead of 8 for speed
                neighbors_values = [
                    gray[i-1, j],     # top
                    gray[i, j+1],     # right  
                    gray[i+1, j],     # bottom
                    gray[i, j-1]      # left
                ]
                
                # Create binary pattern
                binary_val = 0
                for k, neighbor in enumerate(neighbors_values):
                    if neighbor > center:
                        binary_val += (1 << k)
                
                lbp[i, j] = binary_val
        
        return lbp

if __name__ == "__main__":
    comparator = ImageComparator('sample.png', 'output.png')
    results = comparator.run_all_comparisons()
    
    print("\n=== Image Comparison Results ===")
    print(f"SSIM Score: {results['ssim']:.4f} (higher is better, 1.0 is perfect)")
    print(f"MSE Score: {results['mse']:.2f} (lower is better, 0 is perfect)")
    print(f"PSNR Score: {results['psnr']:.2f} dB (higher is better)")
    print(f"Feature Match Confidence: {results['match_confidence']:.4f} (higher is better)")
    
    print("\n=== Edge Quality Analysis ===")
    print(f"Basic Edge Similarity: {results['edge_similarity']:.4f} (higher is better)")
    print(f"Fine Detail Edge Preservation: {results['edge_similarity_tight']:.4f} (higher is better)")
    print(f"Coarse Edge Preservation: {results['edge_similarity_loose']:.4f} (higher is better)")
    print(f"Edge Density Preservation: {results['edge_density_score']:.4f} (higher is better)")
    print(f"Edge Continuity Score: {results['edge_continuity_score']:.4f} (higher is better)")
    print(f"Overall Edge Quality: {results['edge_quality_score']:.4f} (higher is better)")
    
    print("\n=== Density Analysis ===")
    print(f"Overall Content Density: {results['overall_content_density']:.4f}")
    print(f"Text Density: {results['text_density']:.4f}")
    print(f"Image Content Density: {results['image_content_density']:.4f}")
    print(f"White Space Ratio: {results['white_space_ratio']:.4f}")
    print(f"Dense Clusters: {results['num_dense_clusters']}")
    
    print("\n=== Rotation Detection ===")
    print(f"Rotation Angle: {results['rotation_angle']:.2f}°")
    print(f"Rotation Confidence: {results['rotation_confidence']:.4f}")
    print(f"Is Rotated: {results['is_rotated']}")
    
    print("\n=== Noise Analysis ===")
    print(f"Noise Level: {results['noise_level']:.2f}")
    print(f"Noise Quality: {results['noise_quality']}")
    print(f"SNR: {results['snr_db']:.2f} dB")
    print(f"Salt & Pepper Ratio: {results['salt_pepper_ratio']:.6f}")
    
    print("\n=== Texture Analysis ===")
    print(f"Texture Complexity: {results['texture_complexity']}")
    print(f"Contrast Measure: {results['contrast_measure']:.2f}")
    print(f"Texture Energy: {results['texture_energy']:.2f}")
    print(f"Roughness: {results['roughness']:.2f}")
    
    print("\n=== Geometric Analysis ===")
    print(f"Aspect Ratio: {results['aspect_ratio']:.3f}")
    print(f"Rectangularity: {results['avg_rectangularity']:.4f}")
    print(f"Horizontal Symmetry: {results['horizontal_symmetry']:.4f}")
    print(f"Vertical Symmetry: {results['vertical_symmetry']:.4f}")
    print(f"Geometric Quality: {results['geometric_quality']}")
    
    print(f"\nColor Histogram Similarity: {results['histogram_similarity']:.4f} (higher is better)")
    print(f"Average Color Difference (Delta-E): {results['avg_delta_e']:.2f} (lower is better)")
    print(f"Significant Difference: {results['significant_diff_percent']:.2f}% of pixels")
    print(f"Overall Quality Score: {results['overall_quality']:.4f} (higher is better)")
    
    print("\n=== Color/Monochrome Detection ===")
    print(f"Color Type Assessment: {results['color_type_assessment']}")
    print(f"Original Image: {'Grayscale' if results['original_is_grayscale'] else 'Color'} (confidence: {results['original_confidence']:.3f})")
    print(f"Processed Image: {'Grayscale' if results['processed_is_grayscale'] else 'Color'} (confidence: {results['processed_confidence']:.3f})")
    print(f"Color Consistency: {'Yes' if results['color_consistency'] else 'No'}")
    print(f"Color Preservation Score: {results['color_preservation_score']:.4f}")
    
    print(f"\nOriginal Image Details:")
    print(f"  Max Channel Difference: {results['original_max_channel_diff']:.2f}")
    print(f"  Mean Saturation: {results['original_mean_saturation']:.2f}")
    print(f"  Color Pixel Ratio: {results['original_color_pixel_ratio']:.4f}")
    print(f"  Channel Correlation: {results['original_min_correlation']:.6f}")
    print(f"  Identical RGB Ratio: {results['original_identical_channel_ratio']:.4f}")
    
    print(f"\nProcessed Image Details:")
    print(f"  Max Channel Difference: {results['processed_max_channel_diff']:.2f}")
    print(f"  Mean Saturation: {results['processed_mean_saturation']:.2f}")
    print(f"  Color Pixel Ratio: {results['processed_color_pixel_ratio']:.4f}")
    print(f"  Channel Correlation: {results['processed_min_correlation']:.6f}")
    print(f"  Identical RGB Ratio: {results['processed_identical_channel_ratio']:.4f}")
    
    print("\nVisualization files generated:")
    print("- feature_matches.png: Matching features between images")
    print("- edges_original.png: Edge detection for original image")
    print("- edges_processed.png: Edge detection for processed image")
    print("- edges_original_tight.png: Fine detail edge detection for original image")
    print("- edges_processed_tight.png: Fine detail edge detection for processed image")
    print("- edges_original_contours.png: Contour visualization for original image")
    print("- edges_processed_contours.png: Contour visualization for processed image")
    print("- visual_diff.png: Heatmap showing differences between images")
    print("- density_heatmap.png: Content density distribution")
    print("- rotation_detection.png: Detected lines for rotation analysis")
    print("- noise_analysis.png: Noise visualization")
    print("- texture_lbp.png: Local Binary Pattern texture analysis")
    print("- texture_contrast.png: Texture contrast visualization")
    print("- density_heatmap.png: Heatmap of image density characteristics")
    print("- rotation_detection.png: Visualization of detected rotation angles")
    print("- noise_analysis.png: Visualization of noise analysis results")
    print("- texture_lbp.png: Visualization of texture analysis (LBP)")
    print("- texture_contrast.png: Visualization of texture analysis (contrast)")