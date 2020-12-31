package chassis

import "github.com/coolbaobei/chassix/apierrors"

//NewAPIError new api error
func NewAPIError(code int, msg, desc string) *apierrors.APIError {
	return apierrors.New(code, msg, desc)
}
