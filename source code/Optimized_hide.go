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

const ( // all in byte

	//FileHeaderSize : standard size of file header
	FileHeaderSize = 14
	//BMPInfoHeaderSize : standard size of bmpinfo header
	BMPInfoHeaderSize = 40
)

//ReadAllFromFile comments:
// Read all bytes from a file
func ReadAllFromFile(path string) []byte {
	all, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
		return []byte{}
	}
	return all

}

//WriteAllToFile : Write all data to a file.
func WriteAllToFile(data []byte, path string) {
	if err := ioutil.WriteFile(path, data, 0666); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
		return
	}
}

//ProduceImg : Output the bmp file through the indepensible three parts.
// @param imp_path. Output path for the bmp image.
// @param fh, bh, PixelArray. File header, bmpinfo header, pixel array.
// @return possible errors for output
func ProduceImg(ImgPath string, fh []byte, bh []byte, PixelArray []byte) error {
	f, err := os.OpenFile(ImgPath, os.O_RDWR|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	f.Write(fh)
	f.Write(bh)
	f.Write(PixelArray)
	if err := f.Close(); err != nil {
		return err
	}
	return nil

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
// @return FileHeader. File heder of 14 bytes.
// @return BMPInfoHeader. Bmpinfo header of 40 bytes.
// @return PixelArray. Pixel array of bytes.
func GetPartsOfBmp(ImgPath string) ([]byte, []byte, []byte) {
	// TODO Your code here
	TotalFile := ReadAllFromFile(ImgPath)
	PixelArrayStart := FileHeaderSize + BMPInfoHeaderSize
	FileHeader := TotalFile[:FileHeaderSize:FileHeaderSize]
	BMPInfoHeader := TotalFile[FileHeaderSize:PixelArrayStart:PixelArrayStart]
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
	BMPWidthSlice := BMPInfoHeader[4:8:8]
	BMPHeightSlice := BMPInfoHeader[8:12:12]
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
	return FileHeader, BMPInfoHeader, PixelArray
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

//HideText : Hide information into the pixel array
// @param HideData. The text to be hidden
// @param PixelArray. The original pixel array
// @return the modified pixel data, which hides text.
func HideText(HideData []byte, PixelArray []byte) []byte {
	// TODO Your code here
	var TextLength int
	TextLength = len(HideData)
	//convert integer to byte slice
	TextLengthByteSlice := make([]byte, 4)
	TextLengthByteSlice[0] = byte(TextLength & 0x000000FF)
	TextLengthByteSlice[1] = byte((TextLength & 0x0000FF00) >> 8)
	TextLengthByteSlice[2] = byte((TextLength & 0x00FF0000) >> 16)
	TextLengthByteSlice[3] = byte((TextLength & 0xFF000000) >> 24)

	var InsertPlace int64
	InsertPlace = 0
	for i := 0; i < 4; i++ {
		PixelArray, InsertPlace = InsertData(TextLengthByteSlice[i], PixelArray, InsertPlace)
	}
	for i := 0; i < TextLength; i++ {
		PixelArray, InsertPlace = InsertData(HideData[i], PixelArray, InsertPlace)
	}
	return PixelArray
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

//ShowText : Restore the hidden text from the pixel array.
// @param PixelArray. Pixel array in bmp file.
// @return. The hidden text in byte array.
func ShowText(PixelArray []byte) []byte {
	// TODO Your code here
	TextLengthByteSlice := make([]byte, 4)
	var RestorePlace int64
	RestorePlace = 0
	for i := 0; i < 4; i++ {
		TextLengthByteSlice[i], RestorePlace = RestoreData(PixelArray, RestorePlace)
	}
	TextLength := _4byte2int(TextLengthByteSlice)
	OutputInfo := make([]byte, TextLength)
	for i := 0; i < TextLength; i++ {
		OutputInfo[i], RestorePlace = RestoreData(PixelArray, RestorePlace)
	}
	return OutputInfo
}

//HideProcedure : The module to hide text into a image.
func HideProcedure(SourceImgPath string, HideFilePath string, DestImgPath string) {
	fmt.Printf("Hide %v into %v -> %v\n", HideFilePath, SourceImgPath, DestImgPath)
	FileHeader, BMPInfoHeader, PixelArray := GetPartsOfBmp(SourceImgPath)
	HideData := ReadAllFromFile(HideFilePath)
	NewPixelArray := HideText(HideData, PixelArray)
	ProduceImg(DestImgPath, FileHeader, BMPInfoHeader, NewPixelArray)
}

//ShowProcedure : The module to show information from a image.
func ShowProcedure(SourceImgPath string, DataPath string) {
	fmt.Printf("Show hidden text from %v, then write it to %v\n",
		SourceImgPath, DataPath)
	_, _, PixelArray := GetPartsOfBmp(SourceImgPath)
	info := ShowText(PixelArray)
	WriteAllToFile(info, DataPath)

}

//PrintTheUsage : Print the usage of the program
func PrintTheUsage() {
	fmt.Fprintf(os.Stderr, "* hide args: hide <SourceImgPath> <HideFilePath> "+
		"<DestImgPath>\n")
	fmt.Fprintf(os.Stderr, "* show args: show <ImgPath> <data_file>\n")
}

func main() {
	// please do not change any of the following code,
	// or do anything to subvert it.
	if len(os.Args) < 2 {
		PrintTheUsage()
		return
	}
	action := os.Args[1]
	switch action {
	case "hide":
		{
			if len(os.Args) < 5 {
				PrintTheUsage()
			} else {
				HideProcedure(os.Args[2], os.Args[3], os.Args[4])
			}
		}
	case "show":
		{
			if len(os.Args) < 4 {
				PrintTheUsage()
			} else {
				ShowProcedure(os.Args[2], os.Args[3])
			}
		}
	default:
		PrintTheUsage()
	}
}
