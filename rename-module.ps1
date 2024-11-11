$NEW_MODULE = "github.com/sherwin-77/golang-todos"
$OLD_MODULE = "github.com/sherwin-77/go-echo-template"

# Update the module name in go.mod
go mod edit -module $NEW_MODULE

# Rename all imported modules in .go files recursively
Get-ChildItem -Path . -Filter '*.go' -Recurse | ForEach-Object {
    $content = Get-Content -Path $_.FullName
    $updatedContent = $content -replace $OLD_MODULE, $NEW_MODULE
    $updatedContent | Set-Content -Path $_.FullName
}
