# Contribuindo para o Gobiru

## Processo de Desenvolvimento

1. Fork o repositório
2. Crie uma branch para sua feature: `git checkout -b feature/nome-da-feature`
3. Faça suas alterações
4. Execute os testes: `go test ./...`
5. Execute o linter: `golangci-lint run`
6. Commit suas mudanças: `git commit -m 'feat: adiciona nova funcionalidade'`
7. Push para sua branch: `git push origin feature/nome-da-feature`
8. Abra um Pull Request

## Padrões de Código

- Siga as convenções do Go: `go fmt`
- Use nomes descritivos para variáveis e funções
- Adicione comentários quando necessário
- Mantenha funções pequenas e focadas
- Escreva testes para novas funcionalidades

## Commits

Seguimos o padrão Conventional Commits:

- `feat:` nova funcionalidade
- `fix:` correção de bug
- `docs:` alterações na documentação
- `style:` formatação, ponto e vírgula, etc
- `refactor:` refatoração de código
- `test:` adição ou correção de testes
- `chore:` alterações em arquivos de build, etc

## Testes

- Todos os PRs devem incluir testes
- Mantenha a cobertura de testes alta
- Use `go test -race` para verificar race conditions

## Documentação

- Atualize o README.md quando necessário
- Documente novas funcionalidades
- Mantenha a documentação do OpenAPI atualizada

## Dúvidas

Abra uma issue para discutir mudanças maiores antes de começar o trabalho. 