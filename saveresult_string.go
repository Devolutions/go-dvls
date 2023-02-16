// Code generated by "stringer -type=SaveResult -trimprefix SaveResult"; DO NOT EDIT.

package dvls

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[SaveResultError-0]
	_ = x[SaveResultSuccess-1]
	_ = x[SaveResultAccessDenied-2]
	_ = x[SaveResultInvalidData-3]
	_ = x[SaveResultAlreadyExists-4]
	_ = x[SaveResultMaximumReached-5]
	_ = x[SaveResultNotFound-6]
	_ = x[SaveResultLicenseExpired-7]
	_ = x[SaveResultUnknown-8]
	_ = x[SaveResultTwoFactorTypeNotConfigured-9]
	_ = x[SaveResultWebApiRedirectToLogin-10]
	_ = x[SaveResultDuplicateLoginEmail-11]
}

const _SaveResult_name = "ErrorSuccessAccessDeniedInvalidDataAlreadyExistsMaximumReachedNotFoundLicenseExpiredUnknownTwoFactorTypeNotConfiguredWebApiRedirectToLoginDuplicateLoginEmail"

var _SaveResult_index = [...]uint8{0, 5, 12, 24, 35, 48, 62, 70, 84, 91, 117, 138, 157}

func (i SaveResult) String() string {
	if i >= SaveResult(len(_SaveResult_index)-1) {
		return "SaveResult(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _SaveResult_name[_SaveResult_index[i]:_SaveResult_index[i+1]]
}
