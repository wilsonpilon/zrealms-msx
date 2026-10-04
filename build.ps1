<#
.SYNOPSIS
    Script oficial de automação de build, testes, versionamento X.Y.Z e empacotamento ZIP do Z-Realm (zrealm-msx).

.DESCRIPTION
    Regras de Versionamento:
    - Z (Build / Patch): Incrementado automaticamente a cada execução deste script.
    - Y (Feature / Minor): Incrementado via parâmetro -BumpFeature.
    - X (Major): Incrementado via parâmetro -BumpMajor.

    Fluxo de Trabalho:
    1. Lê e atualiza a versão semântica X.Y.Z.
    2. Roda a suite de testes unitários (go test ./...).
    3. Compila os binários executáveis em dist/bin/.
    4. Copia os manuais e licenças para dist/.
    5. Empacota tudo em um arquivo ZIP em dist/ pronto para GitHub Releases.

.PARAMETER BumpFeature
    Incrementa o número Minor/Feature (Y) e reseta Z para 0.

.PARAMETER BumpMajor
    Incrementa o número Major (X) e reseta Y e Z para 0.

.PARAMETER SkipTests
    Pula a execução dos testes automatizados antes do build.

.EXAMPLE
    ./build.ps1
    Incrementa Z (ex: 0.1.0 -> 0.1.1), testa, compila e gera dist/zrealm-msx-v0.1.1-windows-amd64.zip.

.EXAMPLE
    ./build.ps1 -BumpFeature
    Incrementa Y (ex: 0.1.1 -> 0.2.0), testa, compila e empacota.
#>

[CmdletBinding()]
param (
    [switch]$BumpFeature,
    [switch]$BumpMajor,
    [switch]$KeepVersion,
    [switch]$SkipTests
)

$ErrorActionPreference = "Stop"

$rootDir = $PSScriptRoot
if (-not $rootDir) {
    $rootDir = Get-Location
}


$versionFile = Join-Path $rootDir "VERSION"
$pkgVersionFile = Join-Path $rootDir "pkg/version/VERSION"

# 1. Leitura da versão atual
if (-not (Test-Path $versionFile)) {
    "0.1.0" | Out-File -FilePath $versionFile -Encoding utf8 -NoNewline
}

$currentVersionStr = (Get-Content $versionFile -Raw).Trim()
$parts = $currentVersionStr.Split('.')

if ($parts.Length -lt 3) {
    $parts = @("0", "1", "0")
}

[int]$major = [int]$parts[0]
[int]$minor = [int]$parts[1]
[int]$patch = [int]$parts[2]

# 2. Incremento da versão X.Y.Z
if ($KeepVersion) {
    Write-Host "[VERSION] Versão mantida pelo usuário: X=$major, Y=$minor, Z=$patch" -ForegroundColor Cyan
} elseif ($BumpMajor) {
    $major++
    $minor = 0
    $patch = 0
    Write-Host "[VERSION] Major Bump aplicado: X=$major, Y=$minor, Z=$patch" -ForegroundColor Cyan
} elseif ($BumpFeature) {
    $minor++
    $patch = 0
    Write-Host "[VERSION] Feature/Minor Bump aplicado: X=$major, Y=$minor, Z=$patch" -ForegroundColor Cyan
} else {
    $patch++
    Write-Host "[VERSION] Patch/Build auto-incrementado: X=$major, Y=$minor, Z=$patch" -ForegroundColor Cyan
}

$newVersion = "$major.$minor.$patch"
Write-Host "==========================================================" -ForegroundColor Green
Write-Host "   Z-REALM (zrealm-msx) - BUILD PIPELINE v$newVersion" -ForegroundColor Green
Write-Host "==========================================================" -ForegroundColor Green

# Atualiza arquivos VERSION (sem BOM)
$utf8NoBom = [System.Text.UTF8Encoding]::new($false)
[System.IO.File]::WriteAllText($versionFile, $newVersion, $utf8NoBom)
if (Test-Path (Split-Path $pkgVersionFile)) {
    [System.IO.File]::WriteAllText($pkgVersionFile, $newVersion, $utf8NoBom)
}

# 3. Execução dos Testes Automatizados
if (-not $SkipTests) {
    Write-Host "`n[TEST] Executando suite de testes unitários..." -ForegroundColor Yellow
    Push-Location $rootDir
    try {
        $env:GOTMPDIR = Join-Path $rootDir "tmp"
        if (-not (Test-Path $env:GOTMPDIR)) {
            New-Item -ItemType Directory -Path $env:GOTMPDIR -Force | Out-Null
        }
        go test -v ./pkg/models ./pkg/project ./pkg/exporter ./pkg/script ./pkg/storage ./pkg/version
        if ($LASTEXITCODE -ne 0) {
            throw "A suite de testes falhou com código de saída $LASTEXITCODE. Build abortado!"
        }
        Write-Host "[TEST] Todos os testes passaram com sucesso!" -ForegroundColor Green
    } finally {
        Pop-Location
    }
} else {
    Write-Host "`n[TEST] Testes ignorados pelo usuário (-SkipTests)." -ForegroundColor DarkYellow
}

# 4. Preparação do diretório dist/
$distDir = Join-Path $rootDir "dist"
$distBinDir = Join-Path $distDir "bin"

if (Test-Path $distDir) {
    Remove-Item -Path $distDir -Recurse -Force
}
New-Item -ItemType Directory -Path $distBinDir -Force | Out-Null

# 5. Compilação dos Binários
Write-Host "`n[BUILD] Compilando binários Go para Windows x64..." -ForegroundColor Yellow
Push-Location $rootDir
try {
    $targetExe = Join-Path $distBinDir "zrealm.exe"
    $ldFlags = "-s -w -X 'github.com/zrealm-msx/zrealm/pkg/version.rawVersion=$newVersion'"
    
    go build -ldflags $ldFlags -o $targetExe ./cmd/zrealm
    if ($LASTEXITCODE -ne 0) {
        throw "Falha ao compilar cmd/zrealm com código de saída $LASTEXITCODE"
    }

    # Assina o executável gerado para conformidade com Windows Smart App Control
    $cert = Get-ChildItem Cert:\CurrentUser\My -CodeSigningCert -ErrorAction SilentlyContinue | Select-Object -First 1
    if ($cert) {
        Set-AuthenticodeSignature -Certificate $cert -FilePath $targetExe | Out-Null
    }

    Write-Host "[BUILD] Executável gerado com sucesso: $targetExe" -ForegroundColor Green
    Copy-Item -Path $targetExe -Destination (Join-Path $rootDir "zrealm.exe") -Force
} finally {
    Pop-Location
}

# 6. Cópia de Documentos e Binários MSX para dist/
Write-Host "`n[DOCS] Copiando documentação para dist/..." -ForegroundColor Yellow
$docsToCopy = @("README.md", "LICENSE", "MANUAL.md", "RELEASE.md", "SPEC.md", "OUTLINE.md", "CHANGELOG.md")
foreach ($doc in $docsToCopy) {
    $docPath = Join-Path $rootDir $doc
    if (Test-Path $docPath) {
        Copy-Item -Path $docPath -Destination $distDir -Force
    }
}

Write-Host "`n[MSX] Copiando binários da Engine MSX 2 e imagem DSK para dist/msx/..." -ForegroundColor Yellow
$msxDistDir = Join-Path $distDir "msx"
New-Item -ItemType Directory -Path $msxDistDir -Force | Out-Null

$msxFiles = @(
    "engine_msx/emul/dsk/DOS2_zrealm.dsk",
    "engine_msx/emul/rom/zrealm_demo.rom",
    "engine_msx/emul/dos2/zrealm.com",
    "engine_msx/emul/dos2/HEADER.BIN",
    "engine_msx/emul/dos2/GAME.DAT"
)
foreach ($mf in $msxFiles) {
    $srcPath = Join-Path $rootDir $mf
    if (Test-Path $srcPath) {
        Copy-Item -Path $srcPath -Destination $msxDistDir -Force
        Write-Host "  -> $(Split-Path $srcPath -Leaf) copiado para dist/msx/" -ForegroundColor DarkGreen
    }
}

Write-Host "`n[TOOLS] Copiando ferramentas utilitárias (msxtar) para dist/tools/..." -ForegroundColor Yellow
$toolsDistDir = Join-Path $distDir "tools"
New-Item -ItemType Directory -Path $toolsDistDir -Force | Out-Null
$msxtarSrc = Join-Path $rootDir "MSXgl/tools/build/msxtar/msxtar.exe"
if (Test-Path $msxtarSrc) {
    Copy-Item -Path $msxtarSrc -Destination $toolsDistDir -Force
    Write-Host "  -> msxtar.exe copiado para dist/tools/" -ForegroundColor DarkGreen
}

# 7. Geração do Pacote ZIP para GitHub Release
Write-Host "`n[PACK] Gerando pacote ZIP para distribuição no GitHub..." -ForegroundColor Yellow
$zipFileName = "zrealm-msx-v$newVersion-windows-amd64.zip"
$zipFilePath = Join-Path $distDir $zipFileName

# Cria o arquivo temporário de empacotamento com o conteúdo de dist (sem incluir o próprio zip)
$stagingDir = Join-Path $rootDir "staging_pack"
if (Test-Path $stagingDir) {
    Remove-Item -Path $stagingDir -Recurse -Force
}
New-Item -ItemType Directory -Path $stagingDir -Force | Out-Null

Copy-Item -Path "$distDir\*" -Destination $stagingDir -Recurse -Force

# Comprime o conteúdo de staging para dist/<zipFileName>
Compress-Archive -Path "$stagingDir\*" -DestinationPath $zipFilePath -Force
Remove-Item -Path $stagingDir -Recurse -Force

$zipSize = (Get-Item $zipFilePath).Length / 1KB
Write-Host "[PACK] Pacote ZIP gerado: $zipFilePath ($([math]::Round($zipSize, 2)) KB)" -ForegroundColor Green

# 8. Relatório Final de Conclusão
Write-Host "`n==========================================================" -ForegroundColor Green
Write-Host "   BUILD CONCLUÍDO COM SUCESSO! VERSÃO: v$newVersion" -ForegroundColor Green
Write-Host "   Pasta de Distribuição: $distDir" -ForegroundColor Green
Write-Host "   Pacote de Release:     $zipFilePath" -ForegroundColor Green
Write-Host "==========================================================" -ForegroundColor Green
