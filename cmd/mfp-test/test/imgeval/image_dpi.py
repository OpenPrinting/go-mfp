import cv2
import numpy as np
from scipy import fftpack
import matplotlib.pyplot as plt

def downsample_for_analysis(img, max_size=1024):
    """Downsample image for faster analysis while preserving characteristics"""
    h, w = img.shape[:2]
    if max(h, w) <= max_size:
        return img
    
    scale = max_size / max(h, w)
    new_h, new_w = int(h * scale), int(w * scale)
    return cv2.resize(img, (new_w, new_h), interpolation=cv2.INTER_AREA)

def analyze_frequency_content(image_path):
    """ Frequency analysis using downsampled image"""
    img = cv2.imread(image_path, cv2.IMREAD_GRAYSCALE)
    
    # Downsample for fast processing
    img_small = downsample_for_analysis(img, 512)
    
    # Take FFT of smaller image
    f_transform = fftpack.fft2(img_small)
    f_shift = fftpack.fftshift(f_transform)
    magnitude = np.abs(f_shift)
    
    # Calculate power spectrum in concentric rings
    center = np.array(magnitude.shape) // 2
    y, x = np.ogrid[:magnitude.shape[0], :magnitude.shape[1]]
    radius = np.sqrt((x - center[1])**2 + (y - center[0])**2)
    
    # fewer radius bins for speed
    max_radius = int(np.min(center))
    num_bins = min(50, max_radius)
    radial_profile = []
    
    for i in range(num_bins):
        r = (i + 1) * max_radius / num_bins
        mask = (radius <= r) & (radius > r - max_radius / num_bins)
        if np.any(mask):
            radial_profile.append(np.mean(magnitude[mask]))
        else:
            radial_profile.append(0)
    
    return np.array(radial_profile)

def estimate_effective_dpi(radial_profile, nominal_dpi=600):
    """ DPI estimation using simpler analysis"""
    if len(radial_profile) == 0:
        return 0
    
    max_val = np.max(radial_profile)
    if max_val == 0:
        return 0
    
    normalized = radial_profile / max_val
    
    # Calculate high frequency energy ratio
    # Higher resolution images have more energy in high frequencies
    mid_point = len(normalized) // 2
    low_freq_energy = np.sum(normalized[:mid_point])
    high_freq_energy = np.sum(normalized[mid_point:])
    
    if low_freq_energy == 0:
        freq_ratio = 0
    else:
        freq_ratio = high_freq_energy / low_freq_energy
    
    # Higher ratio indicates better preservation of high frequencies
    if freq_ratio > 0.5:
        return 600
    elif freq_ratio > 0.3:
        return 300
    else:
        return 150

def analyze_sharpness(image_path):
    """ sharpness analysis using sample patches"""
    img = cv2.imread(image_path, cv2.IMREAD_GRAYSCALE)
    
    # Downsample for faster processing
    img_small = downsample_for_analysis(img, 1024)
    
    # Use Laplacian variance as a simple sharpness measure
    laplacian = cv2.Laplacian(img_small, cv2.CV_64F)
    sharpness = laplacian.var()
    
    return sharpness

def analyze_texture_complexity(image_path):
    """Analyze texture complexity - upscaled images have smoother textures"""
    img = cv2.imread(image_path, cv2.IMREAD_GRAYSCALE)
    img_small = downsample_for_analysis(img, 512)
    
    img_float = img_small.astype(np.float32)
    
    kernel_size = 5
    pad = kernel_size // 2
    padded = np.pad(img_float, pad, mode='reflect')
    
    local_stds = []
    for i in range(0, img_float.shape[0], 10):
        for j in range(0, img_float.shape[1], 10):
            window = padded[i:i+kernel_size, j:j+kernel_size]
            local_stds.append(np.std(window))
    
    return np.mean(local_stds) if local_stds else 0

def detect_interpolation_artifacts(image_path):
    """Detect interpolation artifacts typical of upscaling"""
    img = cv2.imread(image_path, cv2.IMREAD_GRAYSCALE)
    img_small = downsample_for_analysis(img, 512)
    
    # Check for smoothness patterns
    # Original images have more random noise, upscaled have smooth gradients
    
    # Calculate gradient magnitude
    grad_x = cv2.Sobel(img_small, cv2.CV_64F, 1, 0, ksize=3)
    grad_y = cv2.Sobel(img_small, cv2.CV_64F, 0, 1, ksize=3)
    gradient_mag = np.sqrt(grad_x**2 + grad_y**2)
    
    grad_hist, _ = np.histogram(gradient_mag, bins=50, range=(0, 100))
    grad_hist = grad_hist / np.sum(grad_hist)
    
    # Calculate entropy of gradient distribution
    grad_entropy = -np.sum(grad_hist * np.log2(grad_hist + 1e-10))
    
    return grad_entropy

def analyze_frequency_content_improved(image_path):
    """Improved frequency analysis with better high-frequency detection"""
    img = cv2.imread(image_path, cv2.IMREAD_GRAYSCALE)
    img_small = downsample_for_analysis(img, 512)
    
    # Apply window function to reduce edge artifacts
    window = np.hanning(img_small.shape[0])[:, None] * np.hanning(img_small.shape[1])
    windowed_img = img_small * window
    
    # Take FFT
    f_transform = np.fft.fft2(windowed_img)
    f_shift = np.fft.fftshift(f_transform)
    magnitude = np.abs(f_shift)
    
    # Calculate radial frequency profile
    center = np.array(magnitude.shape) // 2
    y, x = np.ogrid[:magnitude.shape[0], :magnitude.shape[1]]
    radius = np.sqrt((x - center[1])**2 + (y - center[0])**2)
    
    max_radius = int(np.min(center))
    quarter1 = max_radius // 4
    quarter2 = max_radius // 2
    quarter3 = 3 * max_radius // 4
    
    # Calculate energy in different frequency bands
    low_freq_mask = radius <= quarter1
    mid_freq_mask = (radius > quarter1) & (radius <= quarter2)
    high_freq_mask = (radius > quarter2) & (radius <= quarter3)
    very_high_freq_mask = radius > quarter3
    
    low_energy = np.sum(magnitude[low_freq_mask])
    mid_energy = np.sum(magnitude[mid_freq_mask])
    high_energy = np.sum(magnitude[high_freq_mask])
    very_high_energy = np.sum(magnitude[very_high_freq_mask])
    
    total_energy = low_energy + mid_energy + high_energy + very_high_energy
    
    if total_energy == 0:
        return 0, 0, 0
    
    # Calculate frequency ratios
    high_freq_ratio = (high_energy + very_high_energy) / total_energy
    mid_freq_ratio = mid_energy / total_energy
    
    return high_freq_ratio, mid_freq_ratio, very_high_energy / total_energy

def analyze_edge_sharpness_advanced(image_path):
    """Advanced edge sharpness analysis using multiple methods"""
    img = cv2.imread(image_path, cv2.IMREAD_GRAYSCALE)
    img_small = downsample_for_analysis(img, 1024)
    
    # 1. Sobel-based edge sharpness
    sobel_x = cv2.Sobel(img_small, cv2.CV_64F, 1, 0, ksize=3)
    sobel_y = cv2.Sobel(img_small, cv2.CV_64F, 0, 1, ksize=3)
    sobel_magnitude = np.sqrt(sobel_x**2 + sobel_y**2)
    
    # 2. Canny edge detection for counting sharp edges
    edges = cv2.Canny(img_small, 50, 150)
    edge_density = np.sum(edges > 0) / edges.size
    
    # 3. Calculate edge transition sharpness
    edge_pixels = np.where(edges > 0)
    transition_sharpness = []
    
    for y, x in zip(edge_pixels[0][:1000], edge_pixels[1][:1000]):
        if 2 <= y < img_small.shape[0]-2 and 2 <= x < img_small.shape[1]-2:
            grad_direction = np.arctan2(sobel_y[y, x], sobel_x[y, x])
            
            cos_dir = np.cos(grad_direction + np.pi/2)
            sin_dir = np.sin(grad_direction + np.pi/2)
            
            profile = []
            for i in range(-2, 3):
                ny = int(y + i * sin_dir)
                nx = int(x + i * cos_dir)
                if 0 <= ny < img_small.shape[0] and 0 <= nx < img_small.shape[1]:
                    profile.append(img_small[ny, nx])
            
            if len(profile) >= 5:
                edge_gradient = np.max(np.diff(profile))
                transition_sharpness.append(edge_gradient)
    
    avg_transition_sharpness = np.mean(transition_sharpness) if transition_sharpness else 0
    
    return {
        'sobel_variance': sobel_magnitude.var(),
        'edge_density': edge_density,
        'transition_sharpness': avg_transition_sharpness
    }

def detect_upscaling_patterns(image_path):
    """Detect specific patterns that indicate upscaling algorithms"""
    img = cv2.imread(image_path, cv2.IMREAD_GRAYSCALE)
    img_small = downsample_for_analysis(img, 512)
    
    # 1. Check for bilinear interpolation patterns
    autocorr = cv2.matchTemplate(img_small, img_small[::2, ::2], cv2.TM_CCOEFF_NORMED)
    bilinear_score = np.max(autocorr)
    
    # 2. Analyze pixel value distribution
    hist, _ = np.histogram(img_small, bins=256, range=(0, 256))
    
    # Look for peaks at intermediate values
    smoothed_hist = cv2.GaussianBlur(hist.reshape(-1, 1), (5, 1), 0).flatten()
    peaks = []
    for i in range(1, len(smoothed_hist)-1):
        if smoothed_hist[i] > smoothed_hist[i-1] and smoothed_hist[i] > smoothed_hist[i+1]:
            peaks.append(i)
    
    intermediate_peaks = [p for p in peaks if 10 < p < 245 and smoothed_hist[p] > np.mean(smoothed_hist)]
    interpolation_artifacts = len(intermediate_peaks) / len(peaks) if peaks else 0
    
    # 3. Local Binary Pattern analysis
    from skimage.feature import local_binary_pattern
    lbp = local_binary_pattern(img_small, 8, 1, method='uniform')
    lbp_variance = np.var(lbp)
    
    return {
        'bilinear_score': bilinear_score,
        'interpolation_artifacts': interpolation_artifacts,
        'lbp_variance': lbp_variance
    }

def analyze_image_resolution(image_path):
    """Ccomplete analysis to determine effective resolution"""
    
    print(f"Analyzing {image_path}...")
    
    img = cv2.imread(image_path, cv2.IMREAD_GRAYSCALE)
    if img is None:
        return None
    
    print(f"  Image size: {img.shape}")
   
    print("  Running frequency analysis...")
    freq_profile = analyze_frequency_content(image_path)
    estimated_dpi = estimate_effective_dpi(freq_profile)
    
    print("  Running sharpness analysis...")
    sharpness = analyze_sharpness(image_path)
    
    print("  Running texture analysis...")
    texture_complexity = analyze_texture_complexity(image_path)
    
    print("  Running interpolation detection...")
    gradient_entropy = detect_interpolation_artifacts(image_path)
    
    high_freq_ratio, mid_freq_ratio, very_high_energy_ratio = analyze_frequency_content_improved(image_path)
    
    edge_sharpness = analyze_edge_sharpness_advanced(image_path)
    
    upscaling_patterns = detect_upscaling_patterns(image_path)
    
    results = {
        'estimated_dpi': estimated_dpi,
        'sharpness_score': sharpness,
        'texture_complexity': texture_complexity,
        'gradient_entropy': gradient_entropy,
        'high_freq_ratio': high_freq_ratio,
        'mid_freq_ratio': mid_freq_ratio,
        'very_high_energy_ratio': very_high_energy_ratio,
        **edge_sharpness,
        **upscaling_patterns
    }
    
    score = 0
    
    # Sharpness score (weight: 50%)
    if sharpness > 4600:  # 600 DPI threshold
        score += 50
    elif sharpness > 4200:  # 300 DPI threshold  
        score += 35
    elif sharpness > 3500:  # 150 DPI threshold
        score += 20
    else:
        score += 5
    
    # Gradient entropy (weight: 30%) - lower is better (less interpolation)
    if gradient_entropy < 1.29:
        score += 30
    elif gradient_entropy < 1.32:
        score += 20
    elif gradient_entropy < 1.37:
        score += 10
    else:
        score += 0
    
    # Texture complexity (weight: 15%)
    if texture_complexity > 7.9:
        score += 15
    elif texture_complexity > 7.8:
        score += 10
    else:
        score += 5
    
    # Frequency analysis (weight: 5%)
    if estimated_dpi >= 500:
        score += 5
    elif estimated_dpi >= 250:
        score += 3
    else:
        score += 1
    
    # Final classification
    if score >= 80:
        classification = "600 DPI (Original)"
    elif score >= 55:
        classification = "300 DPI (Effective)" 
    else:
        classification = "150 DPI (Effective)"
    
    results['classification'] = classification
    results['confidence_score'] = score
    return results

def analyze_image_resolution_enhanced(image_path):
    """Enhanced complete analysis with improved methods"""
    
    print(f"Analyzing {image_path}...")
    
    img = cv2.imread(image_path, cv2.IMREAD_GRAYSCALE)
    if img is None:
        return None
    
    print(f"  Image size: {img.shape}")
    
    print("  Running enhanced frequency analysis...")
    high_freq_ratio, mid_freq_ratio, very_high_energy_ratio = analyze_frequency_content_improved(image_path)
    
    print("  Running advanced edge analysis...")
    edge_metrics = analyze_edge_sharpness_advanced(image_path)
    
    print("  Running upscaling pattern detection...")
    upscaling_metrics = detect_upscaling_patterns(image_path)
    
    print("  Running basic metrics...")
    sharpness = analyze_sharpness(image_path)
    texture_complexity = analyze_texture_complexity(image_path)
    gradient_entropy = detect_interpolation_artifacts(image_path)
    
    results = {
        'sharpness_score': sharpness,
        'texture_complexity': texture_complexity,
        'gradient_entropy': gradient_entropy,
        'high_freq_ratio': high_freq_ratio,
        'mid_freq_ratio': mid_freq_ratio,
        'very_high_energy_ratio': very_high_energy_ratio,
        **edge_metrics,
        **upscaling_metrics
    }
 
    score = 0
    
    # 1. Sharpness metrics (40% weight)
    if sharpness > 4600:
        score += 25
    elif sharpness > 4200:
        score += 18
    elif sharpness > 3500:
        score += 10
    else:
        score += 2
    
    # Edge transition sharpness
    if edge_metrics['transition_sharpness'] > 50:
        score += 15
    elif edge_metrics['transition_sharpness'] > 30:
        score += 10
    else:
        score += 5
    
    # 2. Frequency content analysis (25% weight)
    if high_freq_ratio > 0.15:
        score += 15
    elif high_freq_ratio > 0.10:
        score += 10
    else:
        score += 3
    
    if very_high_energy_ratio > 0.05:
        score += 10
    elif very_high_energy_ratio > 0.02:
        score += 6
    else:
        score += 1
    
    # 3. Interpolation artifact detection (20% weight)
    if gradient_entropy < 1.29:
        score += 12
    elif gradient_entropy < 1.32:
        score += 8
    elif gradient_entropy < 1.37:
        score += 4
    else:
        score += 0
    
    # Bilinear interpolation detection (lower is better)
    if upscaling_metrics['bilinear_score'] < 0.8:
        score += 8
    elif upscaling_metrics['bilinear_score'] < 0.9:
        score += 4
    else:
        score += 0
    
    # 4. Texture and pattern analysis (15% weight)
    if texture_complexity > 7.9:
        score += 8
    elif texture_complexity > 7.8:
        score += 5
    else:
        score += 2
    
    # Local Binary Pattern variance (higher indicates more complex patterns)
    if upscaling_metrics['lbp_variance'] > 15:
        score += 7
    elif upscaling_metrics['lbp_variance'] > 10:
        score += 4
    else:
        score += 1
    
    if score >= 85:
        classification = "600 DPI (Original)"
        confidence = "High"
    elif score >= 65:
        classification = "300 DPI (Effective)"
        confidence = "Medium-High"
    elif score >= 45:
        classification = "150 DPI (Effective)"
        confidence = "Medium"
    else:
        classification = "150 DPI (Effective)"
        confidence = "Low"
    
    results.update({
        'classification': classification,
        'confidence_score': score,
        'confidence_level': confidence
    })
    
    return results

if __name__ == "__main__":
    images = ['600dpi.png', '300dpi.png', '150dpi.png']
    
    print("Enhanced Image Resolution Analysis")
    print("=" * 60)
    
    for image in images:
        print(f"\n{'='*60}")
        results = analyze_image_resolution_enhanced(image)
        if results:
            print(f"\nRESULTS FOR {image}:")
            print(f"Classification: {results['classification']}")
            print(f"Confidence Score: {results['confidence_score']}/100 ({results['confidence_level']})")
            print(f"\nDetailed Metrics:")
            print(f"  • Sharpness Score: {results['sharpness_score']:.1f}")
            print(f"  • Edge Transition Sharpness: {results['transition_sharpness']:.1f}")
            print(f"  • High Frequency Ratio: {results['high_freq_ratio']:.3f}")
            print(f"  • Very High Frequency Energy: {results['very_high_energy_ratio']:.3f}")
            print(f"  • Gradient Entropy: {results['gradient_entropy']:.3f}")
            print(f"  • Texture Complexity: {results['texture_complexity']:.2f}")
            print(f"  • LBP Variance: {results['lbp_variance']:.1f}")
            print(f"  • Bilinear Score: {results['bilinear_score']:.3f}")
            print(f"  • Edge Density: {results['edge_density']:.4f}")
        else:
            print(f"❌ Could not load {image}")
        print(f"{'='*60}")