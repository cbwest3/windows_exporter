// Copyright 2013 The win Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright (c) 2010 The win Authors. All rights reserved.

// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
// 3. The names of the authors may not be used to endorse or promote products
//    derived from this software without specific prior written permission.

// THIS SOFTWARE IS PROVIDED BY THE AUTHORS ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

//go:build windows

package pdh

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

// PDH error codes, which can be returned by all pdh.* functions. Taken from mingw-w64 pdhmsg.h
const (
	CSTATUS_VALID_DATA                     = 0x00000000 // The returned data is valid.
	CSTATUS_NEW_DATA                       = 0x00000001 // The return data value is valid and different from the last sample.
	CSTATUS_NO_MACHINE                     = 0x800007D0 // Unable to connect to the specified computer, or the computer is offline.
	CSTATUS_NO_INSTANCE                    = 0x800007D1
	MORE_DATA                              = 0x800007D2 // The GetFormattedCounterArray* function can return this if there's 'more data to be displayed'.
	CSTATUS_ITEM_NOT_VALIDATED             = 0x800007D3
	RETRY                                  = 0x800007D4
	NO_DATA                                = 0x800007D5 // The query does not currently contain any counters (for example, limited access)
	CALC_NEGATIVE_DENOMINATOR              = 0x800007D6
	CALC_NEGATIVE_TIMEBASE                 = 0x800007D7
	CALC_NEGATIVE_VALUE                    = 0x800007D8
	DIALOG_CANCELLED                       = 0x800007D9
	END_OF_LOG_FILE                        = 0x800007DA
	ASYNC_QUERY_TIMEOUT                    = 0x800007DB
	CANNOT_SET_DEFAULT_REALTIME_DATASOURCE = 0x800007DC
	CSTATUS_NO_OBJECT                      = 0xC0000BB8
	CSTATUS_NO_COUNTER                     = 0xC0000BB9 // The specified counter could not be found.
	CSTATUS_INVALID_DATA                   = 0xC0000BBA // The counter was successfully found, but the data returned is not valid.
	MEMORY_ALLOCATION_FAILURE              = 0xC0000BBB
	INVALID_HANDLE                         = 0xC0000BBC
	INVALID_ARGUMENT                       = 0xC0000BBD // Required argument is missing or incorrect.
	FUNCTION_NOT_FOUND                     = 0xC0000BBE
	CSTATUS_NO_COUNTERNAME                 = 0xC0000BBF
	CSTATUS_BAD_COUNTERNAME                = 0xC0000BC0 // Unable to parse the counter path. Check the format and syntax of the specified path.
	INVALID_BUFFER                         = 0xC0000BC1
	INSUFFICIENT_BUFFER                    = 0xC0000BC2
	CANNOT_CONNECT_MACHINE                 = 0xC0000BC3
	INVALID_PATH                           = 0xC0000BC4
	INVALID_INSTANCE                       = 0xC0000BC5
	INVALID_DATA                           = 0xC0000BC6 // specified counter does not contain valid data or a successful status code.
	NO_DIALOG_DATA                         = 0xC0000BC7
	CANNOT_READ_NAME_STRINGS               = 0xC0000BC8
	LOG_FILE_CREATE_ERROR                  = 0xC0000BC9
	LOG_FILE_OPEN_ERROR                    = 0xC0000BCA
	LOG_TYPE_NOT_FOUND                     = 0xC0000BCB
	NO_MORE_DATA                           = 0xC0000BCC
	ENTRY_NOT_IN_LOG_FILE                  = 0xC0000BCD
	DATA_SOURCE_IS_LOG_FILE                = 0xC0000BCE
	DATA_SOURCE_IS_REAL_TIME               = 0xC0000BCF
	UNABLE_READ_LOG_HEADER                 = 0xC0000BD0
	FILE_NOT_FOUND                         = 0xC0000BD1
	FILE_ALREADY_EXISTS                    = 0xC0000BD2
	NOT_IMPLEMENTED                        = 0xC0000BD3
	STRING_NOT_FOUND                       = 0xC0000BD4
	UNABLE_MAP_NAME_FILES                  = 0x80000BD5
	UNKNOWN_LOG_FORMAT                     = 0xC0000BD6
	UNKNOWN_LOGSVC_COMMAND                 = 0xC0000BD7
	LOGSVC_QUERY_NOT_FOUND                 = 0xC0000BD8
	LOGSVC_NOT_OPENED                      = 0xC0000BD9
	WBEM_ERROR                             = 0xC0000BDA
	ACCESS_DENIED                          = 0xC0000BDB
	LOG_FILE_TOO_SMALL                     = 0xC0000BDC
	INVALID_DATASOURCE                     = 0xC0000BDD
	INVALID_SQLDB                          = 0xC0000BDE
	NO_COUNTERS                            = 0xC0000BDF
	SQL_ALLOC_FAILED                       = 0xC0000BE0
	SQL_ALLOCCON_FAILED                    = 0xC0000BE1
	SQL_EXEC_DIRECT_FAILED                 = 0xC0000BE2
	SQL_FETCH_FAILED                       = 0xC0000BE3
	SQL_ROWCOUNT_FAILED                    = 0xC0000BE4
	SQL_MORE_RESULTS_FAILED                = 0xC0000BE5
	SQL_CONNECT_FAILED                     = 0xC0000BE6
	SQL_BIND_FAILED                        = 0xC0000BE7
	CANNOT_CONNECT_WMI_SERVER              = 0xC0000BE8
	PLA_COLLECTION_ALREADY_RUNNING         = 0xC0000BE9
	PLA_ERROR_SCHEDULE_OVERLAP             = 0xC0000BEA
	PLA_COLLECTION_NOT_FOUND               = 0xC0000BEB
	PLA_ERROR_SCHEDULE_ELAPSED             = 0xC0000BEC
	PLA_ERROR_NOSTART                      = 0xC0000BED
	PLA_ERROR_ALREADY_EXISTS               = 0xC0000BEE
	PLA_ERROR_TYPE_MISMATCH                = 0xC0000BEF
	PLA_ERROR_FILEPATH                     = 0xC0000BF0
	PLA_SERVICE_ERROR                      = 0xC0000BF1
	PLA_VALIDATION_ERROR                   = 0xC0000BF2
	PLA_VALIDATION_WARNING                 = 0x80000BF3
	PLA_ERROR_NAME_TOO_LONG                = 0xC0000BF4
	INVALID_SQL_LOG_FORMAT                 = 0xC0000BF5
	COUNTER_ALREADY_IN_QUERY               = 0xC0000BF6
	BINARY_LOG_CORRUPT                     = 0xC0000BF7
	LOG_SAMPLE_TOO_SMALL                   = 0xC0000BF8
	OS_LATER_VERSION                       = 0xC0000BF9
	OS_EARLIER_VERSION                     = 0xC0000BFA
	INCORRECT_APPEND_TIME                  = 0xC0000BFB
	UNMATCHED_APPEND_COUNTER               = 0xC0000BFC
	SQL_ALTER_DETAIL_FAILED                = 0xC0000BFD
	QUERY_PERF_DATA_TIMEOUT                = 0xC0000BFE
)

// Formatting options for GetFormattedCounterValue().
const (
	FMT_RAW              = 0x00000010
	FMT_ANSI             = 0x00000020
	FMT_UNICODE          = 0x00000040
	FMT_LONG             = 0x00000100 // Return data as a long int.
	FMT_DOUBLE           = 0x00000200 // Return data as a double precision floating point real.
	FMT_LARGE            = 0x00000400 // Return data as a 64 bit integer.
	FMT_NOSCALE          = 0x00001000 // can be OR-ed: Do not apply the counter's default scaling factor.
	FMT_1000             = 0x00002000 // can be OR-ed: multiply the actual value by 1,000.
	FMT_NODATA           = 0x00004000 // can be OR-ed: unknown what this is for, MSDN says nothing.
	FMT_NOCAP100         = 0x00008000 // can be OR-ed: do not cap values > 100.
	PERF_DETAIL_COSTLY   = 0x00010000
	PERF_DETAIL_STANDARD = 0x0000FFFF
)

type (
	HQUERY   HANDLE // query handle
	HCOUNTER HANDLE // counter handle
)

// For struct details, see https://learn.microsoft.com/en-us/windows/win32/api/pdh/ns-pdh-pdh_counter_info_w.
type COUNTER_INFO struct {
	DwLength        uint32
	DwType          uint32
	CVersion        uint32
	CStatus         uint32
	LScale          int32
	LDefaultScale   int32
	DwUserData      *uint32
	DwQueryUserData *uint32
	SzFullPath      *uint16 // pointer to a string
	CounterPath     COUNTER_PATH_ELEMENTS
	SzExplainText   *uint16 // pointer to a string
	DataBuffer      *string
}

// For struct details, see https://learn.microsoft.com/en-us/windows/win32/api/pdh/ns-pdh-pdh_counter_path_elements_w.
type COUNTER_PATH_ELEMENTS struct {
	SzMachineName    *uint16 // pointer to a string
	SzObjectName     *uint16 // pointer to a string
	SzInstanceName   *uint16 // pointer to a string
	SzParentInstance *uint16 // pointer to a string
	DwInstanceIndex  *uint32
	SzCounterName    *uint16 // pointer to a string
}

// Union specialization for double values
type FMT_COUNTERVALUE_DOUBLE struct {
	CStatus     uint32
	DoubleValue float64
}

// Union specialization for 64 bit integer values
type FMT_COUNTERVALUE_LARGE struct {
	CStatus    uint32
	LargeValue int64
}

// Union specialization for long values
type FMT_COUNTERVALUE_LONG struct {
	CStatus   uint32
	LongValue int32
	padding   [4]byte
}

// Union specialization for double values, used by GetFormattedCounterArrayDouble()
type FMT_COUNTERVALUE_ITEM_DOUBLE struct {
	SzName   *uint16 // pointer to a string
	FmtValue FMT_COUNTERVALUE_DOUBLE
}

// Union specialization for 'large' values, used by GetFormattedCounterArrayLarge()
type FMT_COUNTERVALUE_ITEM_LARGE struct {
	SzName   *uint16 // pointer to a string
	FmtValue FMT_COUNTERVALUE_LARGE
}

// Union specialization for long values, used by GetFormattedCounterArrayLong()
type FMT_COUNTERVALUE_ITEM_LONG struct {
	SzName   *uint16 // pointer to a string
	FmtValue FMT_COUNTERVALUE_LONG
}

var (
	// Library
	libpdhDll *windows.LazyDLL

	// Functions
	pdh_AddCounterW               *windows.LazyProc
	pdh_AddEnglishCounterW        *windows.LazyProc
	pdh_CloseQuery                *windows.LazyProc
	pdh_CollectQueryData          *windows.LazyProc
	pdh_ExpandWildCardPath        *windows.LazyProc
	pdh_GetCounterInfo            *windows.LazyProc
	pdh_GetFormattedCounterValue  *windows.LazyProc
	pdh_GetFormattedCounterArrayW *windows.LazyProc
	pdh_OpenQuery                 *windows.LazyProc
	pdh_ValidatePathW             *windows.LazyProc
)

func init() {
	// Library
	libpdhDll = windows.NewLazySystemDLL("pdh.dll")

	// Functions
	pdh_AddCounterW = libpdhDll.NewProc("PdhAddCounterW")
	pdh_AddEnglishCounterW = libpdhDll.NewProc("PdhAddEnglishCounterW")
	pdh_CloseQuery = libpdhDll.NewProc("PdhCloseQuery")
	pdh_CollectQueryData = libpdhDll.NewProc("PdhCollectQueryData")
	pdh_ExpandWildCardPath = libpdhDll.NewProc("PdhExpandWildCardPathW")
	pdh_GetCounterInfo = libpdhDll.NewProc("PdhGetCounterInfoW")
	pdh_GetFormattedCounterValue = libpdhDll.NewProc("PdhGetFormattedCounterValue")
	pdh_GetFormattedCounterArrayW = libpdhDll.NewProc("PdhGetFormattedCounterArrayW")
	pdh_OpenQuery = libpdhDll.NewProc("PdhOpenQuery")
	pdh_ValidatePathW = libpdhDll.NewProc("PdhValidatePathW")
}

// Adds the specified counter to the query. This is the internationalized version. Preferably, use the
// function AddEnglishCounter instead. hQuery is the query handle, which has been fetched by OpenQuery.
// szFullCounterPath is a full, internationalized counter path (this will differ per Windows language version).
// dwUserData is a 'user-defined value', which becomes part of the counter information. To retrieve this value
// later, call GetCounterInfo() and access dwQueryUserData of the COUNTER_INFO structure.
//
// Examples of szFullCounterPath (in an English version of Windows):
//
//	\\Processor(_Total)\\% Idle Time
//	\\Processor(_Total)\\% Processor Time
//	\\LogicalDisk(C:)\% Free Space
//
// To view all (internationalized...) counters on a system, there are three non-programmatic ways: perfmon utility,
// the typeperf command, and the the registry editor. perfmon.exe is perhaps the easiest way, because it's basically a
// full implemention of the pdh.dll API, except with a GUI and all that. The registry setting also provides an
// interface to the available counters, and can be found at the following key:
//
//	HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Perflib\CurrentLanguage
//
// This registry key contains several values as follows:
//
//	1
//	1847
//	2
//	System
//	4
//	Memory
//	6
//	% Processor Time
//	... many, many more
//
// Somehow, these numeric values can be used as szFullCounterPath too:
//
//	\2\6 will correspond to \\System\% Processor Time
//
// The typeperf command may also be pretty easy. To find all performance counters, simply execute:
//
//	typeperf -qx
func AddCounter(hQuery HQUERY, szFullCounterPath string, dwUserData uintptr, phCounter *HCOUNTER) uint32 {
	ptxt, _ := syscall.UTF16PtrFromString(szFullCounterPath)
	ret, _, _ := pdh_AddCounterW.Call(
		uintptr(hQuery),
		uintptr(unsafe.Pointer(ptxt)),
		dwUserData,
		uintptr(unsafe.Pointer(phCounter)))

	return uint32(ret)
}

// Adds the specified language-neutral counter to the query. See the AddCounter function. This function only exists on
// Windows versions higher than Vista.
func AddEnglishCounter(hQuery HQUERY, szFullCounterPath string, dwUserData uintptr, phCounter *HCOUNTER) uint32 {
	if pdh_AddEnglishCounterW.Find() != nil {
		return ERROR_INVALID_FUNCTION
	}

	ptxt, _ := syscall.UTF16PtrFromString(szFullCounterPath)
	ret, _, _ := pdh_AddEnglishCounterW.Call(
		uintptr(hQuery),
		uintptr(unsafe.Pointer(ptxt)),
		dwUserData,
		uintptr(unsafe.Pointer(phCounter)))

	return uint32(ret)
}

// Closes all counters contained in the specified query, closes all handles related to the query,
// and frees all memory associated with the query.
func CloseQuery(hQuery HQUERY) uint32 {
	ret, _, _ := pdh_CloseQuery.Call(uintptr(hQuery))

	return uint32(ret)
}

// Collects the current raw data value for all counters in the specified query and updates the status
// code of each counter. With some counters, this function needs to be repeatedly called before the value
// of the counter can be extracted with GetFormattedCounterValue(). For example, the following code
// requires at least two calls:
//
//	var handle pdh.HQUERY
//	var counterHandle pdh.HCOUNTER
//	ret := pdh.OpenQuery(0, 0, &handle)
//	ret = pdh.AddEnglishCounter(handle, "\\Processor(_Total)\\% Idle Time", 0, &counterHandle)
//	var derp pdh.FMT_COUNTERVALUE_DOUBLE
//
//	ret = pdh.CollectQueryData(handle)
//	fmt.Printf("Collect return code is %x\n", ret) // return code will be CSTATUS_INVALID_DATA
//	ret = pdh.GetFormattedCounterValueDouble(counterHandle, 0, &derp)
//
//	ret = pdh.CollectQueryData(handle)
//	fmt.Printf("Collect return code is %x\n", ret) // return code will be ERROR_SUCCESS
//	ret = pdh.GetFormattedCounterValueDouble(counterHandle, 0, &derp)
//
// The CollectQueryData will return an error in the first call because it needs two values for
// displaying the correct data for the processor idle time. The second call will have a 0 return code.
func CollectQueryData(hQuery HQUERY) uint32 {
	ret, _, _ := pdh_CollectQueryData.Call(uintptr(hQuery))

	return uint32(ret)
}

// Examines the specified computer or log file and returns those counter paths that match the given counter path which contains wildcard characters.
// For more information, see https://learn.microsoft.com/en-us/windows/win32/api/pdh/nf-pdh-pdhexpandwildcardpathw.
func ExpandWildCardPath(szDataSource *uint16, szWildCardPath *uint16, mszExpandedPathList *uint16, pcchPathListLength *uint32, dwFlags *uint32) uint32 {
	ret, _, _ := pdh_ExpandWildCardPath.Call(
		uintptr(unsafe.Pointer(szDataSource)),
		uintptr(unsafe.Pointer(szWildCardPath)),
		uintptr(unsafe.Pointer(mszExpandedPathList)),
		uintptr(unsafe.Pointer(pcchPathListLength)),
		uintptr(unsafe.Pointer(dwFlags)))

	return uint32(ret)
}

// Retrieves information about a counter, such as data size, counter type, path, and user-supplied data values.
// For more information, see https://learn.microsoft.com/en-us/windows/win32/api/pdh/nf-pdh-pdhgetcounterinfow.
func GetCounterInfo(hCounter HCOUNTER, bRetrieveExplainText uintptr, pdwBufferSize *uint32, lpBuffer *COUNTER_INFO) uint32 {
	ret, _, _ := pdh_GetCounterInfo.Call(
		uintptr(hCounter),
		bRetrieveExplainText,
		uintptr(unsafe.Pointer(pdwBufferSize)),
		uintptr(unsafe.Pointer(lpBuffer)))

	return uint32(ret)
}

// Formats the given hCounter using a 'double'. The result is set into the specialized union struct pValue.
// This function does not directly translate to a Windows counterpart due to union specialization tricks.
func GetFormattedCounterValueDouble(hCounter HCOUNTER, lpdwType *uint32, pValue *FMT_COUNTERVALUE_DOUBLE) uint32 {
	ret, _, _ := pdh_GetFormattedCounterValue.Call(
		uintptr(hCounter),
		uintptr(FMT_DOUBLE),
		uintptr(unsafe.Pointer(lpdwType)),
		uintptr(unsafe.Pointer(pValue)))

	return uint32(ret)
}

// Formats the given hCounter using a large int (int64). The result is set into the specialized union struct pValue.
// This function does not directly translate to a Windows counterpart due to union specialization tricks.
func GetFormattedCounterValueLarge(hCounter HCOUNTER, lpdwType *uint32, pValue *FMT_COUNTERVALUE_LARGE) uint32 {
	ret, _, _ := pdh_GetFormattedCounterValue.Call(
		uintptr(hCounter),
		uintptr(FMT_LARGE),
		uintptr(unsafe.Pointer(lpdwType)),
		uintptr(unsafe.Pointer(pValue)))

	return uint32(ret)
}

// Formats the given hCounter using a 'long'. The result is set into the specialized union struct pValue.
// This function does not directly translate to a Windows counterpart due to union specialization tricks.
//
// BUG(krpors): Testing this function on multiple systems yielded inconsistent results. For instance,
// the pValue.LongValue kept the value '192' on test system A, but on B this was '0', while the padding
// bytes of the struct got the correct value. Until someone can figure out this behaviour, prefer to use
// the Double or Large counterparts instead. These functions provide actually the same data, except in
// a different, working format.
func GetFormattedCounterValueLong(hCounter HCOUNTER, lpdwType *uint32, pValue *FMT_COUNTERVALUE_LONG) uint32 {
	ret, _, _ := pdh_GetFormattedCounterValue.Call(
		uintptr(hCounter),
		uintptr(FMT_LONG),
		uintptr(unsafe.Pointer(lpdwType)),
		uintptr(unsafe.Pointer(pValue)))

	return uint32(ret)
}

// Returns an array of formatted counter values. Use this function when you want to format the counter values of a
// counter that contains a wildcard character for the instance name. The itemBuffer must a slice of type FMT_COUNTERVALUE_ITEM_DOUBLE.
// An example of how this function can be used:
//
//	okPath := "\\Process(*)\\% Processor Time" // notice the wildcard * character
//
//	// ommitted all necessary stuff ...
//
//	var bufSize uint32
//	var bufCount uint32
//	var size uint32 = uint32(unsafe.Sizeof(pdh.FMT_COUNTERVALUE_ITEM_DOUBLE{}))
//	var emptyBuf [1]pdh.FMT_COUNTERVALUE_ITEM_DOUBLE // need at least 1 addressable null ptr.
//
//	for {
//		// collect
//		ret := pdh.CollectQueryData(queryHandle)
//		if ret == pdh.ERROR_SUCCESS {
//			ret = pdh.GetFormattedCounterArrayDouble(counterHandle, &bufSize, &bufCount, &emptyBuf[0]) // uses null ptr here according to MSDN.
//			if ret == pdh.MORE_DATA {
//				filledBuf := make([]pdh.FMT_COUNTERVALUE_ITEM_DOUBLE, bufCount*size)
//				ret = pdh.GetFormattedCounterArrayDouble(counterHandle, &bufSize, &bufCount, &filledBuf[0])
//				for i := 0; i < int(bufCount); i++ {
//					c := filledBuf[i]
//					var s string = pdh.UTF16PtrToString(c.SzName)
//					fmt.Printf("Index %d -> %s, value %v\n", i, s, c.FmtValue.DoubleValue)
//				}
//
//				filledBuf = nil
//				// Need to at least set bufSize to zero, because if not, the function will not
//				// return MORE_DATA and will not set the bufSize.
//				bufCount = 0
//				bufSize = 0
//			}
//
//			time.Sleep(2000 * time.Millisecond)
//		}
//	}
func GetFormattedCounterArrayDouble(hCounter HCOUNTER, lpdwBufferSize *uint32, lpdwBufferCount *uint32, itemBuffer *FMT_COUNTERVALUE_ITEM_DOUBLE) uint32 {
	ret, _, _ := pdh_GetFormattedCounterArrayW.Call(
		uintptr(hCounter),
		uintptr(FMT_DOUBLE),
		uintptr(unsafe.Pointer(lpdwBufferSize)),
		uintptr(unsafe.Pointer(lpdwBufferCount)),
		uintptr(unsafe.Pointer(itemBuffer)))

	return uint32(ret)
}

// Returns an array of formatted counter values. Use this function when you want to format the counter values of a
// counter that contains a wildcard character for the instance name. The itemBuffer must a slice of type FMT_COUNTERVALUE_ITEM_LARGE.
// For an example usage, see GetFormattedCounterArrayDouble.
func GetFormattedCounterArrayLarge(hCounter HCOUNTER, lpdwBufferSize *uint32, lpdwBufferCount *uint32, itemBuffer *FMT_COUNTERVALUE_ITEM_LARGE) uint32 {
	ret, _, _ := pdh_GetFormattedCounterArrayW.Call(
		uintptr(hCounter),
		uintptr(FMT_LARGE),
		uintptr(unsafe.Pointer(lpdwBufferSize)),
		uintptr(unsafe.Pointer(lpdwBufferCount)),
		uintptr(unsafe.Pointer(itemBuffer)))

	return uint32(ret)
}

// Returns an array of formatted counter values. Use this function when you want to format the counter values of a
// counter that contains a wildcard character for the instance name. The itemBuffer must a slice of type FMT_COUNTERVALUE_ITEM_LONG.
// For an example usage, see GetFormattedCounterArrayDouble.
//
// BUG(krpors): See description of GetFormattedCounterValueLong().
func GetFormattedCounterArrayLong(hCounter HCOUNTER, lpdwBufferSize *uint32, lpdwBufferCount *uint32, itemBuffer *FMT_COUNTERVALUE_ITEM_LONG) uint32 {
	ret, _, _ := pdh_GetFormattedCounterArrayW.Call(
		uintptr(hCounter),
		uintptr(FMT_LONG),
		uintptr(unsafe.Pointer(lpdwBufferSize)),
		uintptr(unsafe.Pointer(lpdwBufferCount)),
		uintptr(unsafe.Pointer(itemBuffer)))

	return uint32(ret)
}

// Creates a new query that is used to manage the collection of performance data.
// szDataSource is a null terminated string that specifies the name of the log file from which to
// retrieve the performance data. If 0, performance data is collected from a real-time data source.
// dwUserData is a user-defined value to associate with this query. To retrieve the user data later,
// call GetCounterInfo and access dwQueryUserData of the COUNTER_INFO structure. phQuery is
// the handle to the query, and must be used in subsequent calls. This function returns a
// constant error code, or ERROR_SUCCESS if the call succeeded.
func OpenQuery(szDataSource uintptr, dwUserData uintptr, phQuery *HQUERY) uint32 {
	ret, _, _ := pdh_OpenQuery.Call(
		szDataSource,
		dwUserData,
		uintptr(unsafe.Pointer(phQuery)))

	return uint32(ret)
}

// Validates a path. Will return ERROR_SUCCESS when ok, or CSTATUS_BAD_COUNTERNAME when the path is
// erroneous.
func ValidatePath(path string) uint32 {
	ptxt, _ := syscall.UTF16PtrFromString(path)
	ret, _, _ := pdh_ValidatePathW.Call(uintptr(unsafe.Pointer(ptxt)))

	return uint32(ret)
}
