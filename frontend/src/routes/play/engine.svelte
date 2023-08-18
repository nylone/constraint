<script>
    import { Table, TableBody, TableBodyCell, TableBodyRow } from 'flowbite-svelte'

    // Props 
    export let lobby_id = ""
    export let nickname = ""

    // Variables
    //let rows = Array.apply(null, Array(8)).map(function () {})
    //let cols = Array.apply(null, Array(8)).map(function () {})
    let grid = Array.from(Array(8), () => new Array(8))
    let id = 0

    // Connects to the socket
    let socket_url = `ws://localhost:8080/?lobby=${lobby_id}&nickname=${nickname}`
    const socket = new WebSocket(socket_url)
    
    // Various functions
    function send_move(x, y) {
        socket.send(JSON.stringify({id: id, pos:{x: x, y: y}}))
        id += 1
        grid[y][x] = 'x'
    }


</script>

<div>
    <Table>
        <TableBody tableBodyClass="border-4">
            {#each grid as row, row_index}
                <TableBodyRow>
                    {#each grid as col, col_index}
                    <TableBodyCell
                        id={`col-${col_index+1}`} 
                        tdClass="border px-6 py-4 whitespace-nowrap font-medium "
                        on:click={() => send_move(col_index, row_index)}
                    >
                    {#if grid[row_index][col_index] != undefined}
                        {grid[row_index][col_index]}
                    {/if}
                    </TableBodyCell>
                    {/each}
                </TableBodyRow>
            {/each}
        </TableBody>
    </Table>
</div>