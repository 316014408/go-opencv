// Copyright 2013 <me@cwchang.me>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opencv

//#include "opencv.h"
//#cgo linux  pkg-config: opencv
//#cgo darwin pkg-config: opencv
//#cgo freebsd pkg-config: opencv
//#cgo windows LDFLAGS: -lopencv_core242.dll -lopencv_imgproc242.dll -lopencv_photo242.dll -lopencv_highgui242.dll -lstdc++
import "C"
import (
	//"errors"
	//"log"
	"unsafe"
)

const (
	CV_INTER_NN       = int(C.CV_INTER_NN)
	CV_INTER_LINEAR   = int(C.CV_INTER_LINEAR)
	CV_INTER_CUBIC    = int(C.CV_INTER_CUBIC)
	CV_INTER_AREA     = int(C.CV_INTER_AREA)
	CV_INTER_LANCZOS4 = int(C.CV_INTER_LANCZOS4)
)

func Resize(src *IplImage, width, height, interpolation int) *IplImage {
	if width == 0 && height == 0 {
		panic("Width and Height cannot be 0 at the same time")
	}
	if width == 0 {
		ratio := float64(height) / float64(src.Height())
		width = int(float64(src.Width()) * ratio)
	} else if height == 0 {
		ratio := float64(width) / float64(src.Width())
		height = int(float64(src.Height()) * ratio)
	}

	dst := CreateImage(width, height, src.Depth(), src.Channels())
	C.cvResize(unsafe.Pointer(src), unsafe.Pointer(dst), C.int(interpolation))
	return dst
}

func Crop(src *IplImage, x, y, width, height int) *IplImage {
	r := C.cvRect(C.int(x), C.int(y), C.int(width), C.int(height))
	rect := Rect(r)

	src.SetROI(rect)
	dest := CreateImage(width, height, src.Depth(), src.Channels())
	Copy(src, dest, nil)
	src.ResetROI()

	return dest
}

func CreateContourType() *ContourType {
	return &ContourType{CV_RETR_EXTERNAL, CV_CHAIN_APPROX_SIMPLE, Point{0, 0}}
}

func (this *ContourType) FindContours(image *IplImage) []*Contour {
	storage := C.cvCreateMemStorage(0)
	header_size := (C.size_t)(unsafe.Sizeof(C.CvContour{}))
	var seq *C.CvSeq
	C.cvFindContours(
		unsafe.Pointer(image),
		storage,
		&seq,
		C.int(header_size),
		this.mode,
		this.method,
		C.cvPoint(C.int(this.offset.X), C.int(this.offset.Y)))

	var contours []*Contour
	for i := 0; i < (int)(seq.total); i++ {
		contour := (*Contour)((*_Ctype_CvContour)(unsafe.Pointer(C.cvGetSeqElem(seq, C.int(i)))))
		contours = append(contours, contour)
	}

	storage_c := (*C.CvMemStorage)(storage)
	C.cvReleaseMemStorage(&storage_c)

	return contours
}
