<script>
    import {Heading, Hr ,P, Span} from 'flowbite-svelte'
    import {Blockquote} from 'flowbite-svelte'
</script>

<div>
    <Heading tag="h1" class="mb-4">Rules</Heading>
    <P size="xl" weight="extrabold" class="mb-3">Constraint is a simple 1v1 game inspired by tic-tac-toe, <Span underline>but with a twist </Span></P>
    <Hr class="my-8"/>
    <P size="xl">The game is played in turns, every player gets 1 turn per round.</P>
    <P size="xl">During their turn, a player has to place its mark following the "constraint":</P>
    <Blockquote border class="p-2 my-4"> you may only place your mark in one of the empty slots surrounding the last mark placed by your opponent</Blockquote>
    <P size="xl">The constraint obviously isn't applicable during the turn of whoever plays first.</P>
    <P size="xl">It's easier to explain how to lose rather than how to win</P>
    <Blockquote border class="p-2 my-4"> A player has lost once no more legal placements of their mark are possible</Blockquote>
</div>


