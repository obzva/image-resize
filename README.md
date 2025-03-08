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
