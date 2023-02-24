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
	"errors"
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/prometheus-community/windows_exporter/log"
	"golang.org/x/sys/windows"
)

// Windows error codes from https://learn.microsoft.com/en-us/windows/win32/debug/system-error-codes--0-499-
const (
	ERROR_SUCCESS          = 0x0 // The operation completed successfully.
	ERROR_INVALID_FUNCTION = 0x1 // Incorrect function.
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

// Mapping of hex codes to Windows PDH error names, from
// https://learn.microsoft.com/en-us/windows/win32/perfctrs/pdh-error-codes.
var Errors = map[uint32]string{
	CSTATUS_VALID_DATA:                     "PDH_CSTATUS_VALID_DATA",
	CSTATUS_NEW_DATA:                       "PDH_CSTATUS_NEW_DATA",
	CSTATUS_NO_MACHINE:                     "PDH_CSTATUS_NO_MACHINE",
	CSTATUS_NO_INSTANCE:                    "PDH_CSTATUS_NO_INSTANCE",
	MORE_DATA:                              "PDH_MORE_DATA",
	CSTATUS_ITEM_NOT_VALIDATED:             "PDH_CSTATUS_ITEM_NOT_VALIDATED",
	RETRY:                                  "PDH_RETRY",
	NO_DATA:                                "PDH_NO_DATA",
	CALC_NEGATIVE_DENOMINATOR:              "PDH_CALC_NEGATIVE_DENOMINATOR",
	CALC_NEGATIVE_TIMEBASE:                 "PDH_CALC_NEGATIVE_TIMEBASE",
	CALC_NEGATIVE_VALUE:                    "PDH_CALC_NEGATIVE_VALUE",
	DIALOG_CANCELLED:                       "PDH_DIALOG_CANCELLED",
	END_OF_LOG_FILE:                        "PDH_END_OF_LOG_FILE",
	ASYNC_QUERY_TIMEOUT:                    "PDH_ASYNC_QUERY_TIMEOUT",
	CANNOT_SET_DEFAULT_REALTIME_DATASOURCE: "PDH_CANNOT_SET_DEFAULT_REALTIME_DATASOURCE",
	CSTATUS_NO_OBJECT:                      "PDH_CSTATUS_NO_OBJECT",
	CSTATUS_NO_COUNTER:                     "PDH_CSTATUS_NO_COUNTER",
	CSTATUS_INVALID_DATA:                   "PDH_CSTATUS_INVALID_DATA",
	MEMORY_ALLOCATION_FAILURE:              "PDH_MEMORY_ALLOCATION_FAILURE",
	INVALID_HANDLE:                         "PDH_INVALID_HANDLE",
	INVALID_ARGUMENT:                       "PDH_INVALID_ARGUMENT",
	FUNCTION_NOT_FOUND:                     "PDH_FUNCTION_NOT_FOUND",
	CSTATUS_NO_COUNTERNAME:                 "PDH_CSTATUS_NO_COUNTERNAME",
	CSTATUS_BAD_COUNTERNAME:                "PDH_CSTATUS_BAD_COUNTERNAME",
	INVALID_BUFFER:                         "PDH_INVALID_BUFFER",
	INSUFFICIENT_BUFFER:                    "PDH_INSUFFICIENT_BUFFER",
	CANNOT_CONNECT_MACHINE:                 "PDH_CANNOT_CONNECT_MACHINE",
	INVALID_PATH:                           "PDH_INVALID_PATH",
	INVALID_INSTANCE:                       "PDH_INVALID_INSTANCE",
	INVALID_DATA:                           "PDH_INVALID_DATA",
	NO_DIALOG_DATA:                         "PDH_NO_DIALOG_DATA",
	CANNOT_READ_NAME_STRINGS:               "PDH_CANNOT_READ_NAME_STRINGS",
	LOG_FILE_CREATE_ERROR:                  "PDH_LOG_FILE_CREATE_ERROR",
	LOG_FILE_OPEN_ERROR:                    "PDH_LOG_FILE_OPEN_ERROR",
	LOG_TYPE_NOT_FOUND:                     "PDH_LOG_TYPE_NOT_FOUND",
	NO_MORE_DATA:                           "PDH_NO_MORE_DATA",
	ENTRY_NOT_IN_LOG_FILE:                  "PDH_ENTRY_NOT_IN_LOG_FILE",
	DATA_SOURCE_IS_LOG_FILE:                "PDH_DATA_SOURCE_IS_LOG_FILE",
	DATA_SOURCE_IS_REAL_TIME:               "PDH_DATA_SOURCE_IS_REAL_TIME",
	UNABLE_READ_LOG_HEADER:                 "PDH_UNABLE_READ_LOG_HEADER",
	FILE_NOT_FOUND:                         "PDH_FILE_NOT_FOUND",
	FILE_ALREADY_EXISTS:                    "PDH_FILE_ALREADY_EXISTS",
	NOT_IMPLEMENTED:                        "PDH_NOT_IMPLEMENTED",
	STRING_NOT_FOUND:                       "PDH_STRING_NOT_FOUND",
	UNABLE_MAP_NAME_FILES:                  "PDH_UNABLE_MAP_NAME_FILES",
	UNKNOWN_LOG_FORMAT:                     "PDH_UNKNOWN_LOG_FORMAT",
	UNKNOWN_LOGSVC_COMMAND:                 "PDH_UNKNOWN_LOGSVC_COMMAND",
	LOGSVC_QUERY_NOT_FOUND:                 "PDH_LOGSVC_QUERY_NOT_FOUND",
	LOGSVC_NOT_OPENED:                      "PDH_LOGSVC_NOT_OPENED",
	WBEM_ERROR:                             "PDH_WBEM_ERROR",
	ACCESS_DENIED:                          "PDH_ACCESS_DENIED",
	LOG_FILE_TOO_SMALL:                     "PDH_LOG_FILE_TOO_SMALL",
	INVALID_DATASOURCE:                     "PDH_INVALID_DATASOURCE",
	INVALID_SQLDB:                          "PDH_INVALID_SQLDB",
	NO_COUNTERS:                            "PDH_NO_COUNTERS",
	SQL_ALLOC_FAILED:                       "PDH_SQL_ALLOC_FAILED",
	SQL_ALLOCCON_FAILED:                    "PDH_SQL_ALLOCCON_FAILED",
	SQL_EXEC_DIRECT_FAILED:                 "PDH_SQL_EXEC_DIRECT_FAILED",
	SQL_FETCH_FAILED:                       "PDH_SQL_FETCH_FAILED",
	SQL_ROWCOUNT_FAILED:                    "PDH_SQL_ROWCOUNT_FAILED",
	SQL_MORE_RESULTS_FAILED:                "PDH_SQL_MORE_RESULTS_FAILED",
	SQL_CONNECT_FAILED:                     "PDH_SQL_CONNECT_FAILED",
	SQL_BIND_FAILED:                        "PDH_SQL_BIND_FAILED",
	CANNOT_CONNECT_WMI_SERVER:              "PDH_CANNOT_CONNECT_WMI_SERVER",
	PLA_COLLECTION_ALREADY_RUNNING:         "PDH_PLA_COLLECTION_ALREADY_RUNNING",
	PLA_ERROR_SCHEDULE_OVERLAP:             "PDH_PLA_ERROR_SCHEDULE_OVERLAP",
	PLA_COLLECTION_NOT_FOUND:               "PDH_PLA_COLLECTION_NOT_FOUND",
	PLA_ERROR_SCHEDULE_ELAPSED:             "PDH_PLA_ERROR_SCHEDULE_ELAPSED",
	PLA_ERROR_NOSTART:                      "PDH_PLA_ERROR_NOSTART",
	PLA_ERROR_ALREADY_EXISTS:               "PDH_PLA_ERROR_ALREADY_EXISTS",
	PLA_ERROR_TYPE_MISMATCH:                "PDH_PLA_ERROR_TYPE_MISMATCH",
	PLA_ERROR_FILEPATH:                     "PDH_PLA_ERROR_FILEPATH",
	PLA_SERVICE_ERROR:                      "PDH_PLA_SERVICE_ERROR",
	PLA_VALIDATION_ERROR:                   "PDH_PLA_VALIDATION_ERROR",
	PLA_VALIDATION_WARNING:                 "PDH_PLA_VALIDATION_WARNING",
	PLA_ERROR_NAME_TOO_LONG:                "PDH_PLA_ERROR_NAME_TOO_LONG",
	INVALID_SQL_LOG_FORMAT:                 "PDH_INVALID_SQL_LOG_FORMAT",
	COUNTER_ALREADY_IN_QUERY:               "PDH_COUNTER_ALREADY_IN_QUERY",
	BINARY_LOG_CORRUPT:                     "PDH_BINARY_LOG_CORRUPT",
	LOG_SAMPLE_TOO_SMALL:                   "PDH_LOG_SAMPLE_TOO_SMALL",
	OS_LATER_VERSION:                       "PDH_OS_LATER_VERSION",
	OS_EARLIER_VERSION:                     "PDH_OS_EARLIER_VERSION",
	INCORRECT_APPEND_TIME:                  "PDH_INCORRECT_APPEND_TIME",
	UNMATCHED_APPEND_COUNTER:               "PDH_UNMATCHED_APPEND_COUNTER",
	SQL_ALTER_DETAIL_FAILED:                "PDH_SQL_ALTER_DETAIL_FAILED",
	QUERY_PERF_DATA_TIMEOUT:                "PDH_QUERY_PERF_DATA_TIMEOUT",
}

type (
	HQUERY   uintptr // query handle
	HCOUNTER uintptr // counter handle
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
	nullPtr *uint16

	// Library
	libpdhDll *windows.LazyDLL

	// Functions
	addCounterW               *windows.LazyProc
	addEnglishCounterW        *windows.LazyProc
	closeQuery                *windows.LazyProc
	collectQueryData          *windows.LazyProc
	expandWildCardPath        *windows.LazyProc
	getCounterInfo            *windows.LazyProc
	getFormattedCounterValue  *windows.LazyProc
	getFormattedCounterArrayW *windows.LazyProc
	openQuery                 *windows.LazyProc
	validatePathW             *windows.LazyProc
)

func init() {
	// Library
	libpdhDll = windows.NewLazySystemDLL("pdh.dll")

	// Functions
	addCounterW = libpdhDll.NewProc("PdhAddCounterW")
	addEnglishCounterW = libpdhDll.NewProc("PdhAddEnglishCounterW")
	closeQuery = libpdhDll.NewProc("PdhCloseQuery")
	collectQueryData = libpdhDll.NewProc("PdhCollectQueryData")
	expandWildCardPath = libpdhDll.NewProc("PdhExpandWildCardPathW")
	getCounterInfo = libpdhDll.NewProc("PdhGetCounterInfoW")
	getFormattedCounterValue = libpdhDll.NewProc("PdhGetFormattedCounterValue")
	getFormattedCounterArrayW = libpdhDll.NewProc("PdhGetFormattedCounterArrayW")
	openQuery = libpdhDll.NewProc("PdhOpenQuery")
	validatePathW = libpdhDll.NewProc("PdhValidatePathW")
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
	ret, _, _ := addCounterW.Call(
		uintptr(hQuery),
		uintptr(unsafe.Pointer(ptxt)),
		dwUserData,
		uintptr(unsafe.Pointer(phCounter)))

	return uint32(ret)
}

// Adds the specified language-neutral counter to the query. See the AddCounter function. This function only exists on
// Windows versions higher than Vista.
func AddEnglishCounter(hQuery HQUERY, szFullCounterPath string, dwUserData uintptr, phCounter *HCOUNTER) uint32 {
	if addEnglishCounterW.Find() != nil {
		return ERROR_INVALID_FUNCTION
	}

	ptxt, _ := syscall.UTF16PtrFromString(szFullCounterPath)
	ret, _, _ := addEnglishCounterW.Call(
		uintptr(hQuery),
		uintptr(unsafe.Pointer(ptxt)),
		dwUserData,
		uintptr(unsafe.Pointer(phCounter)))

	return uint32(ret)
}

// Closes all counters contained in the specified query, closes all handles related to the query,
// and frees all memory associated with the query.
func CloseQuery(hQuery HQUERY) uint32 {
	ret, _, _ := closeQuery.Call(uintptr(hQuery))

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
	ret, _, _ := collectQueryData.Call(uintptr(hQuery))

	return uint32(ret)
}

// Examines the specified computer or log file and returns those counter paths that match the given counter path which contains wildcard characters.
// For more information, see https://learn.microsoft.com/en-us/windows/win32/api/pdh/nf-pdh-pdhexpandwildcardpathw.
func ExpandWildCardPath(szDataSource *uint16, szWildCardPath *uint16, mszExpandedPathList *uint16, pcchPathListLength *uint32, dwFlags *uint32) uint32 {
	ret, _, _ := expandWildCardPath.Call(
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
	ret, _, _ := getCounterInfo.Call(
		uintptr(hCounter),
		bRetrieveExplainText,
		uintptr(unsafe.Pointer(pdwBufferSize)),
		uintptr(unsafe.Pointer(lpBuffer)))

	return uint32(ret)
}

// Formats the given hCounter using a 'double'. The result is set into the specialized union struct pValue.
// This function does not directly translate to a Windows counterpart due to union specialization tricks.
func GetFormattedCounterValueDouble(hCounter HCOUNTER, lpdwType *uint32, pValue *FMT_COUNTERVALUE_DOUBLE) uint32 {
	ret, _, _ := getFormattedCounterValue.Call(
		uintptr(hCounter),
		uintptr(FMT_DOUBLE),
		uintptr(unsafe.Pointer(lpdwType)),
		uintptr(unsafe.Pointer(pValue)))

	return uint32(ret)
}

// Formats the given hCounter using a large int (int64). The result is set into the specialized union struct pValue.
// This function does not directly translate to a Windows counterpart due to union specialization tricks.
func GetFormattedCounterValueLarge(hCounter HCOUNTER, lpdwType *uint32, pValue *FMT_COUNTERVALUE_LARGE) uint32 {
	ret, _, _ := getFormattedCounterValue.Call(
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
	ret, _, _ := getFormattedCounterValue.Call(
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
	ret, _, _ := getFormattedCounterArrayW.Call(
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
	ret, _, _ := getFormattedCounterArrayW.Call(
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
	ret, _, _ := getFormattedCounterArrayW.Call(
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
	ret, _, _ := openQuery.Call(
		szDataSource,
		dwUserData,
		uintptr(unsafe.Pointer(phQuery)))

	return uint32(ret)
}

// Validates a path. Will return ERROR_SUCCESS when ok, or CSTATUS_BAD_COUNTERNAME when the path is
// erroneous.
func ValidatePath(path string) uint32 {
	ptxt, _ := syscall.UTF16PtrFromString(path)
	ret, _, _ := validatePathW.Call(uintptr(unsafe.Pointer(ptxt)))

	return uint32(ret)
}

// TODO (cbwest): Do proper error handling.
func LocalizeAndExpandCounter(path string) (paths []string, instances []string, err error) {

	// Open query if it doesn't exist.
	if QueryHandle == 0 {
		log.Debug("Attempting to open PDH query.")

		// var queryHandle HQUERY
		if ret := OpenQuery(0, userData, &QueryHandle); ret != ERROR_SUCCESS {
			return paths, instances, errors.New(fmt.Sprintf("Failed to open PDH query, %s (0x%X).", Errors[ret], ret))
		}
		// client.QueryHandle = &queryHandle
		log.Debug("Opened PDH query successfully.")
	}

	var counterHandle HCOUNTER
	var ret = AddEnglishCounter(QueryHandle, path, 0, &counterHandle)
	if ret != CSTATUS_VALID_DATA { // Error checking
		return paths, instances, errors.New(fmt.Sprintf("AddEnglishCounter returned %s (0x%X)", Errors[ret], ret))
	}

	// Call GetCounterInfo twice to get buffer size, per
	// https://learn.microsoft.com/en-us/windows/win32/api/pdh/nf-pdh-pdhgetcounterinfoa#remarks.
	var bufSize uint32 = 0
	var retrieveExplainText uint32 = 0
	ret = GetCounterInfo(counterHandle, uintptr(retrieveExplainText), &bufSize, nil)
	if ret != MORE_DATA { // error checking
		return paths, instances, errors.New(fmt.Sprintf("First GetCounterInfo returned %s (0x%X)", Errors[ret], ret))
	}

	var counterInfo COUNTER_INFO
	ret = GetCounterInfo(counterHandle, uintptr(retrieveExplainText), &bufSize, &counterInfo)
	if ret != CSTATUS_VALID_DATA { // error checking
		return paths, instances, errors.New(fmt.Sprintf("Second GetCounterInfo returned %s (0x%X)", Errors[ret], ret))
	}

	// Call ExpandWildCardPath twice, per
	// https://learn.microsoft.com/en-us/windows/win32/api/pdh/nf-pdh-pdhexpandwildcardpathw.
	var flags uint32 = 0
	var pathListLength uint32 = 0
	ret = ExpandWildCardPath(nullPtr, counterInfo.SzFullPath, nullPtr, &pathListLength, &flags)
	if ret != MORE_DATA { // error checking
		return paths, instances, errors.New(fmt.Sprintf("ERROR: First ExpandWildCardPath returned %s (0x%X)", Errors[ret], ret))
	}
	if pathListLength < 1 {
		return paths, instances, errors.New(fmt.Sprintf("pathListLength < 1, is %d", pathListLength))
	}

	// TODO (cbwest): Handle PDH_MORE_DATA from https://learn.microsoft.com/en-us/windows/win32/api/pdh/nf-pdh-pdhexpandwildcardpathw.

	expandedPathList := make([]uint16, pathListLength)
	ret = ExpandWildCardPath(nullPtr, counterInfo.SzFullPath, &expandedPathList[0], &pathListLength, &flags)
	if ret != CSTATUS_VALID_DATA { // error checking
		return paths, instances, errors.New(fmt.Sprintf("Second ExpandWildCardPath returned %s (0x%X)", Errors[ret], ret))
	}

	var expandedPath string = ""
	for i := 0; i < int(pathListLength); i += len(expandedPath) + 1 {
		expandedPath = windows.UTF16PtrToString(&expandedPathList[i])
		if len(expandedPath) < 1 { // expandedPathList has two nulls at the end.
			continue
		}

		// Parse PDH instance from the expanded counter path.
		instanceStartIndex := strings.Index(expandedPath, "(")
		instanceEndIndex := strings.Index(expandedPath, ")")
		if instanceStartIndex < 0 || instanceEndIndex < 0 {
			log.Errorf("Unable to parse PDH counter instance from '%s'", expandedPath)
			continue
		}
		instance := expandedPath[instanceStartIndex+1 : instanceEndIndex]

		if instance == "_Total" { // Skip the _Total instance. That is for users to compute.
			log.Debugf("Skipping instance '_Total' for path '%s'", expandedPath)
			continue
		}
		paths = append(paths, expandedPath)
		instances = append(instances, instance)
	}
	return paths, instances, nil
}

var (
	QueryHandle HQUERY
	userData    uintptr // TODO (cbwest): Figure out where to put this.
)


func CollectQueryData2() error {
	ret := CollectQueryData(QueryHandle)
	if ret != CSTATUS_VALID_DATA { // Error checking
		errors.New(fmt.Sprintf("%s (0x%X)\n", Errors[ret], ret))
	}
	return nil
}

func AddCounter2(path string) (*HCOUNTER, error) {
	var counterHandle HCOUNTER
	// Open query if it doesn't exist.
	if QueryHandle == 0 {
		log.Debug("Attempting to open PDH query.")
		if ret := OpenQuery(0, 0, &QueryHandle); ret != 0 {
			return &counterHandle, errors.New(fmt.Sprintf("Failed to open PDH query, %s (0x%X).", Errors[ret], ret))
		}
		log.Debug("Opened PDH query successfully.")
	}

	ret := AddCounter(QueryHandle, path, userData, &counterHandle)
	if ret != CSTATUS_VALID_DATA {
		return &counterHandle, errors.New(fmt.Sprintf("'%s': %s (0x%X)\n", path, Errors[ret], ret))
	}
	return &counterHandle, nil
}
