# Push to spectra-action
cd "C:\Users\Harshal Patel\Desktop\spectra\spectra-action"
if (Test-Path .git) { Remove-Item -Recurse -Force .git }
git init
git config user.name "harshal patel"
git config user.email "hp842484@gmail.com"
git add .
git commit -m "feat: initial release of Spectra GitHub Action"
git branch -M main
git remote add origin https://github.com/HarshalPatel1972/spectra-action.git
git push -u origin main --force

# Setup and push to homebrew-tap
$tapDir = "C:\Users\Harshal Patel\Desktop\spectra\homebrew-tap"
if (-not (Test-Path $tapDir)) { mkdir $tapDir }
cd $tapDir
if (Test-Path .git) { Remove-Item -Recurse -Force .git }
git init
git config user.name "harshal patel"
git config user.email "hp842484@gmail.com"
Set-Content -Path README.md -Value "# Spectra Homebrew Tap`n`nThis repository contains the Homebrew formula for Spectra. It is automatically updated by Goreleaser during the release process."
git add README.md
git commit -m "docs: initialize homebrew tap repository"
git branch -M main
git remote add origin https://github.com/HarshalPatel1972/homebrew-tap.git
git push -u origin main --force
