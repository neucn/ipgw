#!/usr/bin/env pwsh
# edited from deno_install

$ErrorActionPreference = 'Stop'

$BinDir = "$Home\.neucn\bin"

$DownloadedZip = "$env:Temp\ipgw.zip"
$TargetPath = "$BinDir\ipgw.exe"
$Target = if ([System.Environment]::Is64BitOperatingSystem) {
    "windows-amd64"
} else {
    "windows-386"
}

[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$DownloadURL = "https://github.com/neucn/ipgw/releases/latest/download/ipgw-${Target}.zip"

if (!(Test-Path $BinDir)) {
    New-Item $BinDir -ItemType Directory | Out-Null
}

Invoke-WebRequest $DownloadURL -OutFile $DownloadedZip -UseBasicParsing

if (Get-Command Expand-Archive -ErrorAction SilentlyContinue) {
    Expand-Archive $DownloadedZip -Destination $BinDir -Force
} else {
    if (Test-Path $TargetPath) {
        Remove-Item $TargetPath
    }
    Add-Type -AssemblyName System.IO.Compression.FileSystem
    [IO.Compression.ZipFile]::ExtractToDirectory($DownloadedZip, $BinDir)
}

Remove-Item $DownloadedZip

$User = [EnvironmentVariableTarget]::User
$Path = [Environment]::GetEnvironmentVariable('Path', $User)
if (!(";$Path;".ToLower() -like "*;$BinDir;*".ToLower())) {
    [Environment]::SetEnvironmentVariable('Path', "$Path;$BinDir", $User)
    $Env:Path += ";$BinDir"
}

Write-Output "ipgw was installed successfully to $TargetPath"
Write-Output "Run 'ipgw --help' to get started"