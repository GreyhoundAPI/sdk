<?php

namespace GreyhoundApi;

/**
 * Thrown for any non-2xx response, or a transport failure.
 */
class ApiException extends \Exception
{
    /** @var int|null HTTP status code (null for transport failures). */
    public $status;
    /** @var string|null Machine-readable error code from the API. */
    public $code;
    /** @var string|null The request id, for support. */
    public $requestId;
    /** @var array The full error object from the API response. */
    public $details;

    public function __construct(string $message, ?int $status = null, ?string $code = null, ?string $requestId = null, array $details = [])
    {
        parent::__construct($message, (int) ($status ?? 0));
        $this->status    = $status;
        $this->code      = $code;
        $this->requestId = $requestId;
        $this->details   = $details;
    }
}
