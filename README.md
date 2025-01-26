# td

td is a To-do TUI app written in golang

## Installation

### macOS
```bash
brew tap voioo/homebrew-tap
brew install td-tui
```

<details>
<summary>Manual Installation</summary>

```bash
# For Apple Silicon Macs:
curl -LO https://github.com/voioo/td/releases/latest/download/td_darwin_arm64.tar.gz
sudo tar xf td_darwin_arm64.tar.gz -C /usr/local/bin td

# For Intel Macs:
curl -LO https://github.com/voioo/td/releases/latest/download/td_darwin_amd64.tar.gz
sudo tar xf td_darwin_amd64.tar.gz -C /usr/local/bin td
```
</details>

### Arch Linux
```bash
yay -S td-tui
```
or
```bash
paru -S td-tui
```
or
```bash
git clone https://aur.archlinux.org/td-tui.git
cd td-tui
makepkg -si
```

## Ubuntu/Debian
```bash
curl -LO https://github.com/voioo/td/releases/latest/download/td_linux_amd64.tar.gz
sudo tar xf td_linux_amd64.tar.gz -C /usr/local/bin td
```

### RHEL/Fedora/CentOS
```bash
curl -LO https://github.com/voioo/td/releases/latest/download/td_linux_amd64.tar.gz
sudo tar xf td_linux_amd64.tar.gz -C /usr/local/bin td
```

You can also check the releases page on Github and download the one you need.

## Usage

Press '?' to view usage

## Acknowledgements

This project is a derivative of [todo-cli](https://github.com/yuzuy/todo-cli), which is developed by [Ren Ogaki (yuzuy)](https://github.com/yuzuy) for the purposes of learning the Go language. The original code is licensed under the MIT License.

## License

This project is released under the BSD Zero Clause License (0BSD). For more details, please refer to the [LICENSE](LICENSE) file.

---