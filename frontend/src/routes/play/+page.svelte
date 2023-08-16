<script>
    import {Heading} from 'flowbite-svelte'
    import {Input, Button} from 'flowbite-svelte'

    let nickname = ""
    let lobby_id = ""

    function start() {
        let socket_url = `ws://localhost:8080/ws?lobby=${lobby_id}&nickname=${nickname}`
        let socket = new WebSocket(socket_url)

        
    }

</script>

<div>
    <Heading tag="h1" class="mb-4 text-center"> Play </Heading>
    <Input id="nickname" size="lg" placeholder="Nickname" class="mb-2" bind:value={nickname}/>
    <Input id="lobby_id" size="lg" placeholder="Lobby id" bind:value={lobby_id}/>
    <center> <Button id="submit" size="lg" class="mt-6" on:click={start}>Play</Button> </center>
</div>