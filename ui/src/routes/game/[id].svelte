<script>
	import Board from '../../components/board/Board.svelte';
	import { page } from '$app/stores';
</script>

<Board id={$page.params.id} />
