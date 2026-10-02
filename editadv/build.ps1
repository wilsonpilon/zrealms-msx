# ==============================================================================
# Script de Build do Editor de Adventures (Go / PC)
# ==============================================================================
# Execução: powershell -ExecutionPolicy Bypass -File build.ps1
# ==============================================================================

$ErrorActionPreference = "Stop"

$RootDir = $PSScriptRoot
if (-not $RootDir) { $RootDir = (Get-Location).Path }

Write-Host "======================================================================" -ForegroundColor Cyan
Write-Host " [Editor de Adventures] Iniciando Build da Ferramenta em Go" -ForegroundColor Cyan
Write-Host "======================================================================" -ForegroundColor Cyan

# 1. Verificar instalação do Go
Write-Host "`n[1/5] Verificando compilador Go..." -ForegroundColor Yellow
try {
    $goVersion = & go version
    Write-Host "  OK: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "  ERRO: Compilador Go não encontrado no PATH! Instale o Go em https://go.dev" -ForegroundColor Red
    exit 1
}

# 2. Baixar e sincronizar dependências
Write-Host "`n[2/5] Baixando dependências (go mod tidy / download)..." -ForegroundColor Yellow
$CompilerDir = Join-Path $RootDir "compiler"
Push-Location $CompilerDir
try {
    & go mod tidy
    & go mod download
    Write-Host "  OK: Dependências sincronizadas com sucesso." -ForegroundColor Green
} finally {
    Pop-Location
}

# 3. Compilar executável do Editor TUI (edadv.exe)
Write-Host "`n[3/5] Compilando executável do Editor TUI (edadv.exe)..." -ForegroundColor Yellow
$DistDir = Join-Path $RootDir "dist"
if (Test-Path $DistDir) {
    Remove-Item $DistDir -Recurse -Force
}
New-Item -ItemType Directory -Path $DistDir | Out-Null

$ExeOutput = Join-Path $DistDir "edadv.exe"
Push-Location $CompilerDir
try {
    # -ldflags="-s -w" remove informações de debug para reduzir tamanho do binário
    & go build -ldflags="-s -w" -o $ExeOutput ./cmd/edadv
    Write-Host "  OK: edadv.exe gerado com sucesso em $ExeOutput" -ForegroundColor Green
} catch {
    Write-Host "  ERRO ao compilar edadv.exe!" -ForegroundColor Red
    exit 1
} finally {
    Pop-Location
}

# 4. Compilar também o compilador CLI de linha de comando (edadvc.exe)
Write-Host "`n[4/5] Compilando ferramenta CLI (edadvc.exe)..." -ForegroundColor Yellow
$CliOutput = Join-Path $DistDir "edadvc.exe"
Push-Location $CompilerDir
try {
    & go build -ldflags="-s -w" -o $CliOutput ./cmd
    Write-Host "  OK: edadvc.exe gerado com sucesso em $CliOutput" -ForegroundColor Green
} catch {
    Write-Host "  ERRO ao compilar edadvc.exe!" -ForegroundColor Red
    exit 1
} finally {
    Pop-Location
}

# 5. Montar estrutura do pacote distribuível em dist/
Write-Host "`n[5/5] Montando pacote em dist/..." -ForegroundColor Yellow

# Copiar pasta de jogos / modelos de história
$DistGames = Join-Path $DistDir "games"
New-Item -ItemType Directory -Path $DistGames | Out-Null
Copy-Item (Join-Path $RootDir "games\*") $DistGames -Recurse -Force
Write-Host "  - Histórias copiadas para dist/games/" -ForegroundColor Gray

# Copiar fontes da engine C (necessários para F5 / compilação das matrizes)
$DistEngine = Join-Path $DistDir "engine\src"
New-Item -ItemType Directory -Path $DistEngine -Force | Out-Null
Copy-Item (Join-Path $RootDir "engine\src\*") $DistEngine -Recurse -Force
Write-Host "  - Fontes da Engine C copiados para dist/engine/src/" -ForegroundColor Gray

# Criar script de inicialização rápida por duplo clique (iniciar.bat)
$BatContent = @"
@echo off
title Editor de Adventures - MSX
cls
edadv.exe -f games\demo.yaml
"@
Set-Content -Path (Join-Path $DistDir "iniciar.bat") -Value $BatContent -Encoding ASCII
Write-Host "  - Criado dist/iniciar.bat" -ForegroundColor Gray

# Criar arquivo de instruções (LEIA-ME.txt)
$ReadmeContent = @"
======================================================================
  EDITOR DE ADVENTURES - MSX Clean-Room Edition
  Desenvolvido em Go (PC) e C / MSXgl (MSX)
======================================================================

COMO EXECUTAR:
  1. Dê dois cliques em 'iniciar.bat' (ou execute 'edadv.exe' no terminal).
  2. O editor abrirá a aventura de demonstração 'games/demo.yaml'.

ATALHOS DE TECLADO (Interface Estilo Turbo Vision):
  [1]           : Aba de Salas e Mapa de Conexões
  [2]           : Aba de Objetos e Byte de Consistência
  [3]           : Aba de Comandos e Instruções Bytecode
  [4]           : Aba de Mensagens
  [5]           : Aba de Build e Compilação
  [Tab]         : Alternar foco entre listas e formulários
  [Enter]       : Editar ou selecionar item
  [F1]          : Ajuda na tela
  [F2]          : Salvar aventura em arquivo YAML
  [F5]          : Compilar matrizes C (game_data.h e game_data.c)
  [F9]          : Construir ROM do MSX via MSXgl (advent.rom)
  [F10]         : Sair do editor

ESTRUTURA DE PASTAS:
  - edadv.exe   : Editor TUI Interativo
  - edadvc.exe  : Compilador de Linha de Comando (CLI)
  - games/      : Histórias em formato YAML ou JSON
  - engine/src/ : Código-fonte da Engine em C para MSX
  - iniciar.bat : Atalho rápido para iniciar o editor

LINHA DE COMANDO:
  edadv.exe -f games/sua_historia.yaml
  edadvc.exe -i games/sua_historia.yaml -o engine/src

======================================================================
"@
Set-Content -Path (Join-Path $DistDir "LEIA-ME.txt") -Value $ReadmeContent -Encoding UTF8
Write-Host "  - Criado dist/LEIA-ME.txt" -ForegroundColor Gray

# Também copia o edadv.exe gerado para a raiz do projeto para conveniência
Copy-Item $ExeOutput (Join-Path $RootDir "edadv.exe") -Force

Write-Host "`n======================================================================" -ForegroundColor Green
Write-Host " BUILD CONCLUÍDO COM SUCESSO!" -ForegroundColor Green
Write-Host " O pacote autocontido está pronto na pasta: dist\" -ForegroundColor Green
Write-Host " Tamanho do executável: $((Get-Item $ExeOutput).Length / 1MB) MB" -ForegroundColor Green
Write-Host "======================================================================" -ForegroundColor Green
