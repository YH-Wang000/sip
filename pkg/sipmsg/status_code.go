package sipmsg

const (
	Trying               = 100
	Ringing              = 180
	CallIsBeingForwarded = 181
	Queued               = 182
	SessionProgress      = 183
)

// OK 2xx
const OK = 200

// 3xx
const (
	MultipleChoices    = 300
	MovedPermanently   = 301
	MovedTemporarily   = 302
	AlternativeService = 380
)

// Request Failure 4xx
const (
	BadRequest                    = 400
	Unauthorized                  = 401
	PaymentRequired               = 402
	Forbidden                     = 403
	NotFound                      = 404
	MethodNotAllowed              = 405
	NotAcceptable                 = 406
	ProxyAuthenticationRequired   = 407
	RequestTimeout                = 408
	Gone                          = 410
	RequestEntityTooLarge         = 413
	RequestURITooLong             = 414
	UnsupportedMediaType          = 415
	UnsupportedURIScheme          = 416
	BadExtension                  = 420
	ExtensionRequired             = 421
	IntervalTooBrief              = 423
	TemporarilyUnavailable        = 480
	CallOrTransactionDoesNotExist = 481
	LoopDetected                  = 482
	TooManyHops                   = 483
	AddressIncomplete             = 484
	Ambiguous                     = 485
	BusyHere                      = 486
	RequestTerminated             = 487
	NotAcceptableHere             = 488
	RequestPending                = 491
	Undecipherable                = 493
)

// Server Failure 5xx
const (
	ServerInternalError = 500
	NotImplemented      = 501
	BadGateway          = 502
	ServiceUnavailable  = 503
	ServerTimeout       = 504
	VersionNotSupported = 505
	MessageTooLarge     = 513
)

// Global Failure 6xx
const (
	GlobalBusyEverywhere       = 600
	GlobalDecline              = 603
	GlobalDoesNotExistAnywhere = 604
	GlobalNotAcceptable        = 606
)
