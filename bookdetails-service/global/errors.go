package global

type E struct {
	Code int32
	Msg  string
}

var (
	SUCCESS = E{200, "success"}

	ERROR_TOO_MANY_CONNECTIONS = E{429,"too many connections"}

	ERROR_PARAMS_ERROR = E{412,"params error"}

	ERROR_RESOURCE_NOT_FOUND = E{100101,"resouece not found"}

)

