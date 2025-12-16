echo "Configuring zsh theme..."
sed -i 's/ZSH_THEME=".*"/ZSH_THEME="amuse"/' $HOME/.zshrc

go mod tidy
