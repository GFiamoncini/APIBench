# API de Consulta de CEP

Este programa em Go realiza consultas de CEP em dois endpoints diferentes: BrasilAPI e ViaCEP. O objetivo é determinar qual API responde mais rapidamente para um determinado CEP e exibir o resultado. 

## Funcionamento

O programa realiza as seguintes ações:

1. **Consulta Paralela:** Envia requisições simultâneas para duas APIs distintas:
   - **BrasilAPI:** Realiza consulta de informações do CEP usando a API do BrasilAPI.
   - **ViaCEP:** Realiza consulta de informações do CEP usando a API do ViaCEP.

2. **Tempo de Resposta:** Calcula o tempo de resposta de cada API e determina qual delas foi mais rápida.

3. **Exibição dos Resultados:** Imprime na tela os dados recebidos da API mais rápida, incluindo CEP, logradouro, bairro, e UF.

## Estruturas de Dados

- `BraCep`: Representa a resposta JSON da API BrasilAPI.
- `ViaCep`: Representa a resposta JSON da API ViaCEP.

## Como Executar

Certifique-se de ter o Go instalado em sua máquina. Para executar o programa, siga os passos abaixo:

1. Clone este repositório para o seu ambiente local.
2. No terminal, navegue até o diretório do projeto.
3. Execute o comando: 

   ```bash
   go run main.go
