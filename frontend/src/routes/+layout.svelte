<script>
	import '../app.postcss';
</script>

<!-- Flex container -->
<div class="flex justify-center items-center w-100 my-5">
	<div class="flex flex-col">
		<slot />
	</div>
</div>

<style>
	@import '../assets/theme.css';
</style>