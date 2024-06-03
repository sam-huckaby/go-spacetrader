package builder

import (
	"fmt"
)

func Document(title string, content string) (string, error) {
	pageContent := fmt.Sprintf(`
		<html>
			<head>
				<title>%s</title>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<script src="https://cdn.tailwindcss.com"></script>
				<script src="/plugins/htmx.min.js"></script>
			</head>
			<body class="bg-gray-700 text-neutral-300">
				%s
			</body>
		</html>`, title, content)

	return pageContent, nil
}

func Layout_Main(content string) (string, error) {
	LayoutContent := fmt.Sprintf(`
		<div class="w-full flex flex-col justify-center items-center">
			<div class="w-full flex flex-row justify-center items-center p-4">
				<a href="/" class="text-4xl text-neutral-300">Space Trader Interface</a>
			</div>
			<div class="max-w-[960px] w-full">
				%s
			</div>
		</div>
	`, content)

	return LayoutContent, nil
}

func Layout_Fragment(content string) (string, error) {
	LayoutContent := fmt.Sprintf(`
		<div class="max-w-[960px] w-full">
		%s
		</div>
	`, content)

	return LayoutContent, nil
}
