<script>
    import { Table, TableBody, TableBodyCell, TableBodyRow } from 'flowbite-svelte'
    import { Heading, P } from 'flowbite-svelte'

    // Props 
    export let lobby_id = ""
    export let nickname = ""

    // Variables
    let grid = Array.from(Array(8), () => new Array(8))
    let ActionId = 0 
    /* 
        Sets the desired action:
            0 - Add position. This is the default action
            1 - Send message
    */

    let players = 1
    let msg = ''
    let error = false
    let mark
    let MatchOver = false

    // Connects to the socket
    let socket_url = `ws://localhost:8080/?lobby=${lobby_id}&nickname=${nickname}`
    const socket = new WebSocket(socket_url)
    

    // Events handling
    socket.onmessage = function(e) {
        // Gets the output signal
        let signal = JSON.parse(e.data)
        let signal_id = signal["id"]

        // Handles the signals
        switch(signal_id) {
            case 0: // OutputController signal
                if(!signal["successful"]) {
                    msg = signal["error"]
                    error = true
                }
                else {
                    msg = ''
                    error = false
                }
                break;
            case 1: // StartingInfo signal
                if(Object.keys(signal["players"]).length > 1) {
                    players = Object.keys(signal["players"]).length
                }
                mark = signal["players"][nickname] // Assigns the associated mark 
                grid = signal["field"]
                break;
            case 2: // ModelUpdate signal
                // updates the grid with the added position
                let x = signal["pos"]["x"]
                let y = signal["pos"]["y"]
                let placed_mark = signal["mark"]
                grid[y][x] = placed_mark
                
                // checks if there's a winner
                if(signal["winner"] === mark) {
                    msg = `${nickname} won the match`
                    MatchOver = true
                }
                break;
            case 3: // NewClientInfo signal
                players++
                mark = signal["mark"]
                msg = `${signal["nickname"]} has connected to the lobby`
                break;
            case 4: // JoinResponse signal
                if(!signal["successful"]) {
                    msg = signal["error"]
                    error = true
                }
                else {
                    msg = 'Successfully connected'
                    error = false
                }
                break;
            case 5: // ClientLeft signal
                MatchOver = signal["shutdown"]
                msg = `${signal["nickname"]} has left the lobby`
                break;
            case 6: // ChatMessage signal
                break;
        }
        console.log(signal)
    }


    // Various functions
    function send_move(x, y) {
        if(ActionId === 0) {
            socket.send(JSON.stringify({id: ActionId, pos:{x: x, y: y}}))
        }
    }

</script>

<div>
    {#if players < 2} 
        <Heading tag="h4" class="mb-4 text-center">
            Waiting for players ...
        </Heading>
    {/if}
    <P> Players: {players} </P>
    <Table>
        <TableBody tableBodyClass="border-4">
            {#each grid as row, row_index}
                <TableBodyRow>
                    {#each grid as col, col_index}
                    <TableBodyCell
                        id={`col-${col_index}`} 
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
    <Heading tag="h5">Console</Heading>
    {#if msg != '' && !error} 
        <P > {msg} </P> 
    {:else if msg != '' && error}
        <P color="text-red-700 dark:text-red-500"> {msg} </P>
    {/if} 

    {#if MatchOver} 
        <Heading tag="h3"> Match is over </Heading>
    {/if}

</div>