<script>
    import { Table, TableBody, TableBodyCell, TableBodyRow } from 'flowbite-svelte'
    import { Heading, P, Input } from 'flowbite-svelte'

    // Props 
    export let lobby_id = ""
    export let nickname = ""

    // Variables
    let grid = Array.from(Array(8), () => new Array(8))

    let players = 1
    let messages = ['Attempting connection']
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
                    messages.length = messages.push(`ERROR: ${signal["error"]}`)
                    error = true
                }
                else {
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
                    messages.length = messages.push(`${nickname} won the match`)
                    MatchOver = true
                }
                break;
            case 3: // NewClientInfo signal
                players++
                mark = signal["mark"]
                messages.length = messages.push(`${signal["nickname"]} has connected to the lobby`)
                break;
            case 4: // JoinResponse signal
                if(!signal["successful"]) {
                    messages.length = messages.push(`ERROR: ${signal["error"]}`)
                    error = true
                }
                else {
                    messages.length = messages.push('Successfully connected')
                    error = false
                }
                break;
            case 5: // ClientLeft signal
                MatchOver = signal["shutdown"]
                messages.length = messages.push(`${signal["nickname"]} has left the lobby`)
                break;
            case 6: // ChatMessage signal
                messages.length = messages.push(`${signal["by"]}: ${signal["msg"]}`)
                break;
        }
    }

    // Chat message handler
    function send_msg(e) {
        let msg = document.getElementById("chat")
        if(e.key === "Enter" && msg.value != '') {
            socket.send(JSON.stringify({id: 1, msg: msg.value}))
            msg.value = ''
        }
    }

    // Various functions
    function send_move(x, y) {
        socket.send(JSON.stringify({id: 0, pos:{x: x, y: y}}))
    }

</script>

<div>
    {#if players < 2} 
        <Heading tag="h4" class="mb-4 text-center">
            Waiting for players ...
        </Heading>
    {/if}

    {#if MatchOver} 
        <Heading tag="h3"> Match is over </Heading>
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
                        on:click={() => send_move(col_index, row_index, 0)}
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
    {#each messages as msg}
        {#if msg != undefined && !error}
            <P> {msg} </P> 
        {:else if msg != undefined && error}
        <P color="text-red-700 dark:text-red-500"> {msg} </P>
        {/if}
    {/each}

    <Input id="chat" class="mt-5" size="lg" placeholder="Message" on:keydown={send_msg} />

</div>