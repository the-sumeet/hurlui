

export const timingsDescription: Record<string, string> = {
    name_lookup: "DNS resolution time - Time spent resolving the domain name to an IP address. Subsequent requests may be faster due to DNS caching.",
    connect: "TCP connection establishment time - Time to establish the initial TCP connection to the server. Shows 0 when connection is reused.",
    app_connect: "SSL/TLS handshake time - Time spent establishing secure connection encryption. Shows 0 when existing secure connection is reused.",
    pre_transfer: "Time until ready to send request - Total time from start until the request can be sent, includes DNS + TCP + TLS setup.",
    start_transfer: "Time to first byte received - Includes server processing time and network latency. Measures how long server takes to start responding.",
    total: "Complete request duration - Total time for the entire HTTP request/response cycle from start to finish.",
    begin_call: "Call start timestamp - ISO 8601 formatted timestamp indicating when the HTTP request was initiated.",
    end_call: "Call end timestamp - ISO 8601 formatted timestamp indicating when the HTTP response was fully received."
};


export const responseTypes = [
    { value: "json", label: "JSON" },
    { value: "xml", label: "XML" },
    { value: "html", label: "HTML" }
];
