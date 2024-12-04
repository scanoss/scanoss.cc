# README


## Local Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Local Development environment with parameters
| Parameter      | Description                                                                 |
|----------------|-----------------------------------------------------------------------------|
| **scan-root**  | Scanned folder (optional - default $WORKDIR).                               |
| **input**      | Path to results.json file of the scanned project (optional - default $WORKDIR/.scanoss/results.json). |
| **config**     | Path to configuration file (optional - default $HOME/.scanoss/scanoss-lui-settings.json). |
| **apiUrl**     | SCANOSS API URL (optional - default: https://api.osskb.org).                |
| **key**        | SCANOSS API Key token (optional - not required for default OSSKB URL).      |

To start the application with specific arguments, use the following command:

### Example
```shell
wails dev -appargs "--input <resultPath>" 
```

Or you can also run the application using make command:
```shell
make run APPARGS="--scan-root <scanRootPath> --input <resultPath>"
```

All the parameters are optional and have default values. If you change the default values, they will be saved in the configuration file.

## Building

Run the following command to copy the necessary assets and build the application.

```shell
make build
```


## Errors on Ubuntu 24.04
Ubuntu 24.04 includes webkit 4.1 and wails is expecting webkit 4.0
`Perhaps you should add the directory containing 'webkit2gtk-4.0.pc'`

### Solution for Production Environments:
To resolve the issue in a production environment, you can create symbolic links from the WebKit 4.1 libraries to the expected WebKit 4.0 filenames:
```bash
sudo ln -sf /usr/lib/x86_64-linux-gnu/libwebkit2gtk-4.1.so.0 /usr/lib/x86_64-linux-gnu/libwebkit2gtk-4.0.so.37 &&
sudo ln -sf /usr/lib/x86_64-linux-gnu/libjavascriptcoregtk-4.1.so.0 /usr/lib/x86_64-linux-gnu/libjavascriptcoregtk-4.0.so.18
````
### Solution for Development Environments:
If you are working in a development environment, Wails provides a tag to support WebKit 4.1 directly. You can use the following command:<br>
`
wails dev -tags webkit2_41
`






