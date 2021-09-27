<svelte:options tag="data-graph"/>



<script>
import { onMount } from "svelte";
let canvas
import 'chartjs-adapter-moment';

import {
  Chart,
  ArcElement,
  LineElement,
  BarElement,
  PointElement,
  BarController,
  BubbleController,
  DoughnutController,
  LineController,
  PieController,
  PolarAreaController,
  RadarController,
  ScatterController,
  CategoryScale,
  LinearScale,
  LogarithmicScale,
  RadialLinearScale,
  TimeScale,
  TimeSeriesScale,
  Decimation,
  Filler,
  Legend,
  Title,
  Tooltip
} from 'chart.js';

Chart.register(
  ArcElement,
  LineElement,
  BarElement,
  PointElement,
  BarController,
  BubbleController,
  DoughnutController,
  LineController,
  PieController,
  PolarAreaController,
  RadarController,
  ScatterController,
  CategoryScale,
  LinearScale,
  LogarithmicScale,
  RadialLinearScale,
  TimeScale,
  TimeSeriesScale,
  Decimation,
  Filler,
  Legend,
  Title,
  Tooltip
);


let tempData

onMount(async ()=> {
    tempData = {}

    tempData = await fetch('/temperatures').then(res => res.json()).then(data => data)

    let ctx = canvas.getContext('2d');

    console.log(tempData.sensors)

    const colors = ["#FF7A00", "#FF9900", "#FFC700"]

    let datasets = tempData.sensors.map((sensor, index) => {
        return {
            borderWidth: 5,
            borderColor: colors[index],
            backgroundColor: colors[index],
            cubicInterpolationMode: 'monotone',
            tension: 0.4,
            label: sensor.name,
            data: sensor.values.map(val => { console.log(val.x); return {x: new Date(val.x), y: val.y}})
        }
    })

    let data = {
        datasets,
    }

    console.log(data)

    var myChart = new Chart(ctx, {
        type: 'line',
        data: data,
        options: {
            animation: {
                duration: 0
            },
            elements: {
                point: {
                    radius: 0
                }
            },
            scales: {
                x: {
                    type: 'time',
                },
                y: {
                    beginAtZero: true
                }
            }
        },
})});
</script>

<canvas bind:this={canvas} class="chart" height="300px"></canvas>

<style>
    .chart {
        height: 300px;
        max-height: 300px;
    }
</style>