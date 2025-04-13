import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '30s', target: 50 }, // Aumenta a carga para 50 usuários ao longo de 30 segundos
        { duration: '1m', target: 50 },  // Mantém 50 usuários por 1 minuto
        { duration: '30s', target: 0 },  // Desce a carga para 0 usuários em 30 segundos
    ],
};

export default function () {
    // A URL base da sua aplicação
    const url = 'http://localhost:8000/healthz'; // Substitua com o endereço correto se for no Kubernetes

    // Requisição HTTP
    const res = http.get(url);

    // Verificação do status da resposta (deve ser 200)
    check(res, {
        'status is 200': (r) => r.status === 200,
    });

    // Sleep para simular usuários de forma mais realista
    sleep(1);
}
