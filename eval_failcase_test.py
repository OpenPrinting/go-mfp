# Copyright (C) 2026 Mohammad Arman (officialmdarman@gmail.com)
# See LICENSE for license terms and conditions
#
# Validates that the image evaluator correctly rejects known-bad image pairs:
#   - 90-degree rotation (content is rotated, not just paper)
#   - heavy crop (>10% of edges cut off, simulating severe scaling artifact)
# Also confirms a good pair (same image) scores near 1.0.

import sys
import os
import tempfile
import cv2
import numpy as np

# Add imgeval directory to path so ImageComparator can be imported directly.
sys.path.insert(0, os.path.join(os.path.dirname(__file__),
    'cmd', 'mfp-test', 'test', 'imgeval'))

from enhanced_comparison import ImageComparator

THRESHOLD = 0.95


def make_test_image(path: str, w: int = 100, h: int = 75) -> None:
    """Generate a simple grayscale test image with shapes and gradients."""
    img = np.ones((h, w, 3), dtype=np.uint8) * 240  # light-gray background
    # Dark rectangle
    cv2.rectangle(img, (10, 10), (60, 50), (30, 30, 30), -1)
    # Lighter rectangle
    cv2.rectangle(img, (65, 15), (90, 60), (100, 100, 100), -1)
    # Some gradient noise
    noise = np.random.randint(0, 20, (h, w, 3), dtype=np.uint8)
    img = np.clip(img.astype(np.int16) - noise, 0, 255).astype(np.uint8)
    cv2.imwrite(path, img)


def score_pair(ref_path: str, cap_path: str) -> float:
    comp = ImageComparator(ref_path, cap_path)
    results = comp.run_all_comparisons()
    return results['overall_quality']


def run_tests():
    with tempfile.TemporaryDirectory() as tmpdir:
        ref  = os.path.join(tmpdir, 'ref.png')
        good = os.path.join(tmpdir, 'good.png')
        rot  = os.path.join(tmpdir, 'rotated.png')
        crop = os.path.join(tmpdir, 'cropped.png')

        make_test_image(ref)
        import shutil
        shutil.copy(ref, good)  # identical copy — should score ~1.0

        # Rotation fail case: rotate reference 90 degrees clockwise.
        ref_img = cv2.imread(ref)
        rotated = cv2.rotate(ref_img, cv2.ROTATE_90_CLOCKWISE)
        cv2.imwrite(rot, rotated)

        # Crop fail case: remove 15% from each edge (simulates severe scaling).
        h, w = ref_img.shape[:2]
        margin_x = int(w * 0.15)
        margin_y = int(h * 0.15)
        cropped = ref_img[margin_y:h - margin_y, margin_x:w - margin_x]
        cv2.imwrite(crop, cropped)

        pairs = [
            ('near-identical (pipeline sim)', good,   True),
            ('90-deg rotation (must FAIL)',   rot,    False),
            ('heavy crop ~15% (must FAIL)',   crop,   False),
        ]

        all_ok = True
        for label, cap, expect_pass in pairs:
            score = score_pair(ref, cap)
            passed = score >= THRESHOLD
            status = 'PASS' if passed else 'FAIL'
            correct = (passed == expect_pass)
            marker = '✓' if correct else '✗ WRONG'
            if not correct and not expect_pass:
                # Fail case incorrectly passed — critical error.
                all_ok = False
            print(f'  {status}  score={score:.4f}  {label}  {marker}')

        print()
        if all_ok:
            print('FAIL-CASE VALIDATION PASSED: evaluator rejects rotation and crop.')
        else:
            print('FAIL-CASE VALIDATION FAILED: a bad image pair incorrectly passed!')
            sys.exit(1)


if __name__ == '__main__':
    np.random.seed(42)
    run_tests()
