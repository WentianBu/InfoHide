/*
--------------------Copyright Information---------------------

				  Copyright Wentian Bu 2017
		   Powered By csintro.ucas.ac.cn && Wentian Bu

The following parts are written by Wentian Bu:
	function _4byte2int
	function GetPartsOfBmp
	function HideText
	function ShowText
	function InsertData
	function RestoreData

Other Parts are provided by csintro.ucas.ac.cn

This Information Hide program is Open Source.


--------------------Copyright Information---------------------
*/
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

//ReadAllFromFile comments:
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

// Odutput the bmp file through the indepensible three parts.
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
	OutputInteger := (int(bs[0])) | (int(bs[1]) << 8) | (int(bs[2]) << 16) | (int(bs[3]) << 24)
	return OutputInteger
}

// GetPartsOfBmp comments:
// Retrieve three parts of the bmp file: file header, bmpinfo header and pixel
// array. Note the bmp file may contain other parts after the pixel array.
// @param imp_path. The bmp file path.
// @return file_header. File heder of 14 bytes.
// @return bmpinfo_header. Bmpinfo header of 40 bytes.
// @return pixel_array. Pixel array of bytes.
func GetPartsOfBmp(img_path string) ([]byte, []byte, []byte) {
	// TODO Your code here
	TotalFile := ReadAllFromFile(img_path)
	PixelArrayStart := FILE_HEADER_SIZE + BMPINFO_HEADER_SIZE
	file_header := TotalFile[:FILE_HEADER_SIZE:FILE_HEADER_SIZE]
	bmpinfo_header := TotalFile[FILE_HEADER_SIZE:PixelArrayStart:PixelArrayStart]
	/* Notes:
	The <Complex mode> part is used to deal with bitmap files
	which has contents behind the pixel array. In this experiment, the pictures
	provided don't have content behind the pixel array.

	If you don't need this, please comment the <Complex mode> part
	and uncomment the <Easy mode> part.
	*/

	//<Complex mode>
	//<Get the number of pixels>

	// Make slices of height and width
	BMPWidthSlice := bmpinfo_header[4:8:8]
	BMPHeightSlice := bmpinfo_header[8:12:12]
	//convert slices to integers
	var BMPWidth, BMPHeight int
	BMPWidth = _4byte2int(BMPWidthSlice)
	BMPHeight = _4byte2int(BMPHeightSlice)
	//calculate numbers of pixels array bytes

	//</Get the number of pixels>

	PixelArrayEnd := PixelArrayStart + BMPHeight*BMPWidth*3
	PixelArray := TotalFile[PixelArrayStart:PixelArrayEnd]

	//</Complex mode>

	/*
		//<Easy mode>

		PixelArray := TotalFile[PixelArrayStart:]

		//</Easy mode>
	*/
	return file_header, bmpinfo_header, PixelArray
}

//InsertData : Insert data into the byte slice
// @param Data. The data(only one byte) to insert
// @param PixelArray. The pixel array to insert data
// @param InsertPlace. The place in the pixel array to insert data at
// @return the inserted pixel array and the next place
func InsertData(Data byte, PixelArray []byte, InsertPlace int64) ([]byte, int64) {
	for i := 0; i < 4; i++ {
		Low2Bit := Data & 0x3
		PixelArray[InsertPlace] = PixelArray[InsertPlace] & 0xFC
		PixelArray[InsertPlace] = PixelArray[InsertPlace] | Low2Bit
		Data = Data >> 2
		InsertPlace++
	}
	return PixelArray, InsertPlace
}

// Hide information into the pixel array
// @param hide_data. The text to be hidden
// @param pixel_array. The original pixel array
// @return the modified pixel data, which hides text.
func HideText(hide_data []byte, pixel_array []byte) []byte {
	// TODO Your code here
	var TextLength int
	TextLength = len(hide_data)
	//convert integer to byte slice
	TextLengthByteSlice := make([]byte, 4)
	TextLengthByteSlice[0] = byte(TextLength & 0x000000FF)
	TextLengthByteSlice[1] = byte((TextLength & 0x0000FF00) >> 8)
	TextLengthByteSlice[2] = byte((TextLength & 0x00FF0000) >> 16)
	TextLengthByteSlice[3] = byte((TextLength & 0xFF000000) >> 24)

	var InsertPlace int64
	InsertPlace = 0
	for i := 0; i < 4; i++ {
		pixel_array, InsertPlace = InsertData(TextLengthByteSlice[i], pixel_array, InsertPlace)
	}
	for i := 0; i < TextLength; i++ {
		pixel_array, InsertPlace = InsertData(hide_data[i], pixel_array, InsertPlace)
	}
	return pixel_array
}

//RestoreData : Restoe data from file
// @param PixelArray. Pixel array from where to restore data
// @param RestorePlace. The place in the pixel array from where to restore data
// @return data hidden in the pixel array and the next restore place.
func RestoreData(PixelArray []byte, RestorePlace int64) (byte, int64) {
	var Data byte
	var i uint
	for i = 0; i < 4; i++ {
		Low2Bit := PixelArray[RestorePlace] & 0x3
		Low2Bit = Low2Bit << (2 * i)
		Data = Data + Low2Bit
		RestorePlace++
	}
	return Data, RestorePlace
}

// Restore the hidden text from the pixel array.
// @param pixel_array. Pixel array in bmp file.
// @return. The hidden text in byte array.
func ShowText(pixel_array []byte) []byte {
	// TODO Your code here
	TextLengthByteSlice := make([]byte, 4)
	var RestorePlace int64
	RestorePlace = 0
	for i := 0; i < 4; i++ {
		TextLengthByteSlice[i], RestorePlace = RestoreData(pixel_array, RestorePlace)
	}
	TextLength := _4byte2int(TextLengthByteSlice)
	OutputInfo := make([]byte, TextLength)
	for i := 0; i < TextLength; i++ {
		OutputInfo[i], RestorePlace = RestoreData(pixel_array, RestorePlace)
	}
	return OutputInfo
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
