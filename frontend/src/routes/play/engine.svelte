<script>
    import { onMount } from 'svelte'
    import { Table, TableBody, TableBodyCell, TableBodyRow } from 'flowbite-svelte'

    export let lobby_id = ""
    export let nickname = ""

    let rows = Array.apply(null, Array(8)).map(function () {})
    let cols = Array.apply(null, Array(8)).map(function () {})

    onMount(() => {
        let socket_url = `ws://localhost:8080/?lobby=${lobby_id}&nickname=${nickname}`
        let socket = new WebSocket(socket_url)
    })

</script>

<div>
    <Table>
        <TableBody tableBodyClass="border-4">
            {#each rows as row}
            <TableBodyRow>
                {#each cols as col}
                <TableBodyCell tdClass="border px-6 py-4 whitespace-nowrap font-medium "></TableBodyCell>
                {/each}
            </TableBodyRow>
            {/each}
        </TableBody>
    </Table>
</div>