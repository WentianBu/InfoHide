package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	// all in byte
	FILE_HEADER_SIZE    = 14 // standard size of file header
	BMPINFO_HEADER_SIZE = 40 // standard size of bmpinfo header
	LENGTH_FIELD_SIZE   = 16 // size of occupancy in bmp for the length of hidden data
	INFO_UNIT_SIZE      = 4  // size of occupancy in bmp for a byte of hidden data
)

// Read all bytes from a file
func ReadAllFromFile(path string) []byte {
	if all, err := ioutil.ReadFile(path); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
		return []byte{}
	} else {
		return all
	}
}

// Write all data to a file.
func WriteAllToFile(data []byte, path string) {
	if err := ioutil.WriteFile(path, data, 0666); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
		return
	}
}

// Output the bmp file through the indepensible three parts.
// @param imp_path. Output path for the bmp image.
// @param fh, bh, pixel_array. File header, bmpinfo header, pixel array.
// @return possible errors for output
func ProduceImg(img_path string, fh []byte, bh []byte, pixel_array []byte) error {
	if f, err := os.OpenFile(img_path, os.O_RDWR|os.O_CREATE, 0660); err != nil {
		return err
	} else {
		f.Write(fh)
		f.Write(bh)
		f.Write(pixel_array)
		if err := f.Close(); err != nil {
			return err
		} else {
			return nil
		}
	}
}

// Transform bytes to an integer in a little-endian way
// @param bs. byte slice
// @return. The integer value transformed by the slice
func _4byte2int(bs []byte) int {
	// TODO Your code here
}

// Retrieve three parts of the bmp file: file header, bmpinfo header and pixel
// array. Note the bmp file may contain other parts after the pixel array.
// @param imp_path. The bmp file path.
// @return file_header. File heder of 14 bytes.
// @return bmpinfo_header. Bmpinfo header of 40 bytes.
// @return pixel_array. Pixel array of bytes.
func GetPartsOfBmp(img_path string) ([]byte, []byte, []byte) {
	var file_header, bmpinfo_header, pixel_array []byte
	// TODO Your code here

	return file_header, bmpinfo_header, pixel_array
}

// Hide information into the pixel array
// @param hide_data. The text to be hidden
// @param pixel_array. The original pixel array
// @return the modified pixel data, which hides text.
func HideText(hide_data []byte, pixel_array []byte) []byte {
	// TODO Your code here
}

// Restore the hidden text from the pixel array.
// @param pixel_array. Pixel array in bmp file.
// @return. The hidden text in byte array.
func ShowText(pixel_array []byte) []byte {
	// TODO Your code here
}

func HideProcedure(src_img_path string, hide_file_path string, dest_img_path string) {
	fmt.Printf("Hide %v into %v -> %v\n", hide_file_path, src_img_path, dest_img_path)
	file_header, bmpinfo_header, pixel_array := GetPartsOfBmp(src_img_path)
	hide_data := ReadAllFromFile(hide_file_path)
	new_pixel_array := HideText(hide_data, pixel_array)
	ProduceImg(dest_img_path, file_header, bmpinfo_header, new_pixel_array)
}

func ShowProcedure(src_img_path string, data_path string) {
	fmt.Printf("Show hidden text from %v, then write it to %v\n",
		src_img_path, data_path)
	_, _, pixel_array := GetPartsOfBmp(src_img_path)
	info := ShowText(pixel_array)
	WriteAllToFile(info, data_path)

}

func _print_usage() {
	fmt.Fprintf(os.Stderr, "* hide args: hide <src_img_path> <hide_file_path> "+
		"<dest_img_path>\n")
	fmt.Fprintf(os.Stderr, "* show args: show <img_path> <data_file>\n")
}

func main() {
	// please do not change any of the following code,
	// or do anything to subvert it.
	if len(os.Args) < 2 {
		_print_usage()
		return
	} else {
		action := os.Args[1]
		switch action {
		case "hide":
			{
				if len(os.Args) < 5 {
					_print_usage()
				} else {
					HideProcedure(os.Args[2], os.Args[3], os.Args[4])
				}
			}
		case "show":
			{
				if len(os.Args) < 4 {
					_print_usage()
				} else {
					ShowProcedure(os.Args[2], os.Args[3])
				}
			}
		default:
			_print_usage()
		}
	}
}
