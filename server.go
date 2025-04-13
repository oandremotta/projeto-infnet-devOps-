package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

var (
	startedAt = time.Now()

	requestsTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "app_requests_total",
		Help: "Número total de requisições recebidas na raiz (/).",
	})

	healthzChecks = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "app_healthz_checks_total",
		Help: "Número de verificações de healthz.",
	})

	requestDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "app_request_duration_seconds",
		Help:    "Duração das requisições na rota / em segundos.",
		Buckets: prometheus.DefBuckets,
	})

	cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "app_cpu_usage_percent",
		Help: "Porcentagem de uso da CPU.",
	})

	memUsage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "app_memory_usage_percent",
		Help: "Porcentagem de uso da memória.",
	})
)

func init() {
	// Registro das métricas no Prometheus
	prometheus.MustRegister(requestsTotal)
	prometheus.MustRegister(healthzChecks)
	prometheus.MustRegister(requestDuration)
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memUsage)
}

func main() {
	// Coletar métricas do sistema em background
	go collectSystemMetrics()

	// Rotas HTTP
	http.HandleFunc("/", Hello)
	http.HandleFunc("/healthz", Healthz)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("Servidor rodando na porta 80...")
	http.ListenAndServe(":80", nil)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	requestsTotal.Inc()

	name := os.Getenv("NAME")
	if name == "" {
		name = "sem nome definido"
	}
	fmt.Fprintf(w, "Hey, eu sou o %s\nProjeto realizado para o curso da Infnet.", name)

	duration := time.Since(start).Seconds()
	requestDuration.Observe(duration)
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	healthzChecks.Inc()

	duration := time.Since(startedAt)
	if duration.Seconds() < 10 || duration.Seconds() > 320 {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Duration: %v", duration.Seconds())))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func collectSystemMetrics() {
	for {
		// CPU
		cpuPercent, err := cpu.Percent(0, false)
		if err == nil && len(cpuPercent) > 0 {
			cpuUsage.Set(cpuPercent[0])
		}

		// Memória
		memStats, err := mem.VirtualMemory()
		if err == nil {
			memUsage.Set(memStats.UsedPercent)
		}

		time.Sleep(10 * time.Second)
	}
}
