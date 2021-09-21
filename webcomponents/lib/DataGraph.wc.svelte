<svelte:options tag="data-graph"/>

<script>
  import Line from "svelte-chartjs/src/Line.svelte"
</script>

<Line 
width={400}
height={400}
data={{
    datasets: [
        {
            label: "Saunan lämpötila",
            data: [
                {x: 1, y: 20},
                {x: 2, y: 28},
                {x: 3, y: 17},
                {x: 4, y: 19},
            ]
        }
    ],
}} />