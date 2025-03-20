> [!IMPORTANT]  
> I built a image-processing package. Visit [gato](https://github.com/obzva/gato).

# image-resize

A Go package for image resizing that supports multiple interpolation methods for images. (**only supports JPEG right now**)

It includes a command-line interface for testing.

## Features

- Resize images using three interpolation methods:
  - Nearest neighbor
  - Bilinear
  - Bicubic
- Command-line interface for easy testing and usage
- Optional concurrency mode for improved performance

## Usage

### Command Line Interface

```bash
git clone https://github.com/obzva/image-resize.git

cd ./image-resize

go run main.go -p input.jpg -w 800 -h 600 -m bilinear -o output.jpg -c true
```

### Parameters

- `-p`: Path to input image (**required**)
- `-w`: Desired width of output image, defaults to keep the ratio of the original image when omitted (**at least one of two, width or height, is required**)
- `-h`: Desired height of output image, defaults to keep the ratio of the original image when omitted (**at least one of two, width or height, is required**)
- `-m`: Interpolation method, defaults to nearestneighbor when omitted (options: nearestneighbor, bilinear, bicubic)
- `-o`: Output filename, defaults to the method name when omitted
- `-c`: Concurrency mode, defaults to true when omitted

### Example

Let's scale sample image up twice

![input image](/assets/images/test-image.jpg)

input image (500 x 300)

![nearest neighbor output image](/assets/images/nearestneighbor.jpg)

output image using nearest neighbor interpolation (1000 x 600)

![bilinear output image](/assets/images/bilinear.jpg)

output image using bilinear interpolation (1000 x 600)

![bicubic output image](/assets/images/bicubic.jpg)

output image using bicubic interpolation (1000 x 600)

## Package Structure

```
image-resize/
├── main.go                    # Entry point for CLI application
├── imageprocessor/
│   └── imageprocessor.go      # Handles file I/O and manages the image processing workflow
└── interpolator/
    └── interpolator.go        # Implements interpolation algorithms (nearestneighbor, bilinear, bicubic)
    └── interpolator_test.go   # Tests nearest-neighbor and bilinear methods
```

## License

The MIT License (MIT)

Copyright (c) 2025 obzva

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
