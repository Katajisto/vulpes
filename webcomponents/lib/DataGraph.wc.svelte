<svelte:options tag="data-graph"/>



<script>
import { onMount } from "svelte";
let canvas

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

onMount(()=> {
    let ctx = canvas.getContext('2d');
    var myChart = new Chart(ctx, {
    type: 'line',
    data: {
        datasets: [{
            label: 'Saunan lämpötila',
            data: [{x: 10, y: 10}, {x: 20, y: 20}, {x: 30, y: 30}, {x: 40, y: 40}, {x: 50, y: 50}],
            backgroundColor: [
                "#707070"
        ],
            borderWidth: 5,
            borderColor: "#707070"
        }]
    },
    options: {
        scales: {
            y: {
                beginAtZero: true
            }
        }
    }
});})

</script>

<canvas bind:this={canvas} height="70"></canvas>
