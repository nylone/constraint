<script>
    import {Heading} from 'flowbite-svelte'
    import {Input, Button} from 'flowbite-svelte'
    import Engine from './engine.svelte';

    let nickname = ""
    let lobby_id = ""
    let play = true

    // TODO: input check

    function start() {
        play = true
    }

</script>

<div>
    {#if !play}
        <Heading tag="h1" class="mb-4 text-center"> Play </Heading>
        <Input id="nickname" size="lg" placeholder="Nickname" class="mb-2" bind:value={nickname}/>
        <Input id="lobby_id" size="lg" placeholder="Lobby id" bind:value={lobby_id}/>
        <center> <Button id="submit" size="lg" class="mt-6" on:click={start}>Play</Button> </center>

    {:else} 
        <Engine lobby_id={lobby_id} nickname={nickname}/>
    {/if}
</div>