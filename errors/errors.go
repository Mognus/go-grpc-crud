package errors

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Problem struct {
	Type     string       `json:"type"`
	Title    string       `json:"title"`
	Status   int          `json:"status"`
	Detail   string       `json:"detail,omitempty"`
	Instance string       `json:"instance,omitempty"`
	Errors   []FieldError `json:"errors,omitempty"`
	Err      error        `json:"-"`
}

func (p *Problem) Error() string {
	if p.Detail != "" {
		return p.Detail
	}
	return p.Title
}

func NotFound(resource string) *Problem {
	return &Problem{
		Type:   "/problems/not-found",
		Title:  "Not Found",
		Status: 404,
		Detail: fmt.Sprintf("%s not found", resource),
	}
}

func BadRequest(detail string) *Problem {
	return &Problem{
		Type:   "/problems/bad-request",
		Title:  "Bad Request",
		Status: 400,
		Detail: detail,
	}
}

func ValidationFailed(errs []FieldError) *Problem {
	return &Problem{
		Type:   "/problems/validation-failed",
		Title:  "Validation Failed",
		Status: 422,
		Errors: errs,
	}
}

func Unauthorized(detail string) *Problem {
	return &Problem{
		Type:   "/problems/unauthorized",
		Title:  "Unauthorized",
		Status: 401,
		Detail: detail,
	}
}

func Forbidden(detail string) *Problem {
	return &Problem{
		Type:   "/problems/forbidden",
		Title:  "Forbidden",
		Status: 403,
		Detail: detail,
	}
}

func Conflict(detail string) *Problem {
	return &Problem{
		Type:   "/problems/conflict",
		Title:  "Conflict",
		Status: 409,
		Detail: detail,
	}
}

func TooManyRequests(detail string) *Problem {
	return &Problem{
		Type:   "/problems/too-many-requests",
		Title:  "Too Many Requests",
		Status: 429,
		Detail: detail,
	}
}

func Internal(err error) *Problem {
	return &Problem{
		Type:   "/problems/internal-error",
		Title:  "Internal Server Error",
		Status: 500,
		Detail: "An unexpected error occurred",
		Err:    err,
	}
}

func GrpcToHTTP(err error) *Problem {
	st, _ := status.FromError(err)
	msg := st.Message()
	switch st.Code() {
	case codes.NotFound:
		return NotFound(msg)
	case codes.InvalidArgument:
		return BadRequest(msg)
	case codes.AlreadyExists:
		return Conflict(msg)
	case codes.Unauthenticated:
		return Unauthorized(msg)
	case codes.PermissionDenied:
		return Forbidden(msg)
	default:
		return Internal(err)
	}
}
