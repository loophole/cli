package httpserver

const (
	// first %s is logoUrl, second %s is error message
	proxyErrorTemplate = `<!DOCTYPE html>
<html lang="en">
	<head>
	<meta charset="utf-8" />
	<title>Loophole is running...</title>
	<style>
		html {
		height: 100%%;
		}
		body {
			max-height: 100%%;
			font-family: system-ui, -apple-system, "Segoe UI", Roboto, Ubuntu,
				Cantarell, "Noto Sans", sans-serif, BlinkMacSystemFont, "Segoe UI",
				Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji",
				"Segoe UI Symbol";
			overflow-y: auto;
		}
		.container {
			text-align: center;
			width: 800px;
			height: fit-content;

			position: absolute;
			top: 0;
			bottom: 0;
			left: 0;
			right: 0;

			margin: auto;
		}
		.error {
			color: #fa383e;
		}
	</style>
	</head>
	<body>
	<div class="container">
		<img
		src="%s"
		width="500px"
		alt="Loophole"
		/>
		<h1>Congratulations, your tunnel is up and running!</h1>
		<p>
		However... it looks like the application you're trying to expose is not
		available.
		<br />
		<br />
		Original error: <em class="error">%s</em>
		<br />
		<br />
		<small>
			This request would normally end up with 502 status code.
			<br />
			If that's what you intended please restart the tunnel with
			<code>--disable-proxy-error-page</code> option
			<br />
			to remove this page and get regular 502 error.
		</small>
		</p>
	</div>
	</body>
</html>
`
)
