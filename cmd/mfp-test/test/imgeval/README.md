# Image Evaluation

This directory contains image quality evaluation code originally written by
**Sanskar Yaduka** as part of Google Summer of Code 2025 for OpenPrinting.

Original repository:
https://github.com/Sanskary2303/OpenPrinting-Image-Evaluation

## Files

- `enhanced_comparison.py` — `ImageComparator` class that compares two images
  across multiple metrics (SSIM, PSNR, color accuracy, edge quality, etc.) and
  returns an overall quality score between 0.0 and 1.0.
- `image_dpi.py` — detects the actual DPI of a scanned or printed image.

## License

2-clause BSD license. See the original repository for details.

## Dependencies

Install Python dependencies with:

```
pip install -r requirements.txt
```
