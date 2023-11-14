# FlexiProxyHub

Welcome to FlexiProxyHub, the forefront of innovative proxy utility creation. Our platform excels in offering advanced solutions for dynamic and complex proxy management scenarios. FlexiProxyHub is designed to simplify and enhance your proxy management requirements. This proxy solution is intuitive and highly adaptable, making it perfect for both development and production environments. FlexiProxyHub specializes in facilitating complex proxy tasks, like the ReportAsync feature, which establishes a websocket for pages with delayed download initiation, without altering the application itself.

## Features

- **Flexible Routing:** Define custom routes with ease.
- **Environment Modes:** Switch between DEVELOPMENT and PRODUCTION modes.
- **Detailed Logging:** Control log verbosity with body size limits and visible headers.
- **Enhanced Header Management:** Customize headers passed to the backend.
- **Simple Configuration:** Easy setup with environment variables.

## Environment Variables

FlexiProxyHub can be configured using the following environment variables:

- `PROXY_CONFIGURATION`: Define proxy routes and behaviors. 
  - Example: `[{"routes": [{"host": "localhost:8080", "path": ["/teste.txt"]}], "mode": 1, "proxy_to": "http://localhost:8081"}]`
- `ENVIRONMENT`: Set the running environment (`DEVELOPMENT` or `PRODUCTION`).
- `LOG_BODY_MAX_SIZE`: Maximum size of the log body. Default is 255.
- `VISIBLE_HEADERS`: Specify visible headers, separated by a comma, or use `*` for all headers.
- `PROXY_CLIENT_HEADERS_TO_BACKEND`: Define which client headers are forwarded to the backend.
- `LISTEN_PORT`: Port for the proxy to listen on. Default is 8080.
- `LISTEN_HOST`: Host for the proxy to listen on. Default is `localhost`.

## Getting Started

Setting up FlexiProxyHub is straightforward:

1. Clone the repository.
2. Set the required environment variables.
3. Run the application.

## Modules

FlexiProxyHub is designed to be modular, allowing for easy extension and customization. Each module in FlexiProxyHub provides a unique set of functionalities, catering to various proxying needs. Below is a guide to the currently implemented modules in FlexiProxyHub, along with a note on the ongoing restructuring for better modularity.

### 1. AsynchronousReport Module - Advanced Asynchronous Handling

The `AsynchronousReport` module is the first and currently the only module implemented in FlexiProxyHub. 

Specialized in managing delayed download scenarios with our unique AsynchronousReport module, providing robust websocket support for uninterrupted file transfers.

#### Key Features:

- **Websocket-Based Communication:** Utilizes websockets to establish a persistent, real-time communication channel between the client and the backend.
- **Asynchronous File Download:** Allows clients to download files asynchronously, reducing the wait time for file generation and transfer.
- **Resilience to Network Issues:** Minimizes the impact of network interruptions by maintaining an active websocket connection, unlike traditional HTTP connections which are more susceptible to network errors.

#### How It Works:

- When a client requests a file download, `AsynchronousReport` sets up a websocket connection instead of initiating an immediate file transfer.
- The backend then processes the file generation request. During this process, the client's connection remains active via the websocket, but not tied up in waiting for a direct response.
- Once the file is ready, the client is notified through the websocket channel, and the download process begins. This method ensures that the client doesn't have to maintain an open HTTP connection, which can be interrupted by network errors, and instead waits for a signal over the resilient websocket connection.

### Example Configuration

To demonstrate the implementation of the `AsynchronousReport` module, consider the following `PROXY_CONFIGURATION`:

```json
PROXY_CONFIGURATION=[{
    "routes": [{"host": "localhost:8080", "path": ["/teste.txt"]}],
    "mode": 1,
    "proxy_to": "http://localhost:8081"
}]
```


In this configuration:

- The routes define a specific route to be handled by FlexiProxyHub, in this case, any requests to localhost:8080/teste.txt.
The mode is set to 1, which likely indicates the environment mode (either DEVELOPMENT or PRODUCTION).
- The proxy_to field specifies that the incoming requests to the defined route should be forwarded to http://localhost:8081.
- This example demonstrates how FlexiProxyHub can be configured to route specific requests to different backend services, making it a flexible solution for various proxying scenarios.



## Contributing

Contributions to FlexiProxyHub are always welcome! Whether it's bug fixes, feature enhancements, or documentation improvements, feel free to fork the repository and submit a pull request.

## License

FlexiProxyHub is open-source software licensed under the [MIT License](LICENSE).

## Contact

For any questions or feedback, please open an issue in the GitHub repository, and we'll get back to you as soon as possible.

Happy Proxying with FlexiProxyHub!
