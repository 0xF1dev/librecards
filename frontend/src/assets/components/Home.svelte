<script>
    import { onMount } from "svelte";
    import { fade } from "svelte/transition";

    // @ts-ignore
    import { GetCards, DeleteCard } from "../../../wailsjs/go/main/App";

    // @ts-ignore
    import ScrollLeft from "../images/scroll-left.png";
    // @ts-ignore
    import ScrollRight from "../images/scroll-right.png";
    // @ts-ignore
    import Delete from "../images/delete.png";
    // @ts-ignore
    import Edit from "../images/edit.png";

    import { _ } from "svelte-i18n";

    let { onNew, onSelect, onDelete, onEdit } = $props();

    let cards = $state([]);
    let count = $state(0);

    let card = $state({
        id: 0,
        title: "Loading...",
        subject: "Please wait",
        path: "",
    });
    let cardIndex = 0;

    onMount(async () => {
        let index = await GetCards();

        cards = index.cards;
        card = cards[0];
        count = index.count;
    });

    async function scrollLeft() {
        if (cardIndex <= 0) {
            cardIndex = cards.length - 1;
        } else {
            cardIndex -= 1;
        }
        card = cards[cardIndex];
    }

    async function scrollRight() {
        if (cardIndex >= cards.length - 1) {
            cardIndex = 0;
        } else {
            cardIndex += 1;
        }
        card = cards[cardIndex];
    }

    function newFunc() {
        onNew();
    }

    function selectCard(id) {
        onSelect(id);
    }

    async function deleteCard(id) {
        let err = await DeleteCard(id, {
            title: $_("dialogs.delete.title"),
            content: $_("dialogs.delete.content"),
            buttons: [
                $_("dialogs.delete.button_yes"),
                $_("dialogs.delete.button_no"),
            ],
            default: $_("dialogs.delete.default"),
            cancel: $_("dialogs.delete.cancel"),
        });
        if (err === 0) {
            onDelete();
        }
    }

    function edit() {
        onEdit(card.id);
    }
</script>

<div class="wrapper">
    <div class="main">
        <p class="text">
            {$_("home.your_flashcards", { values: { count: count } })}
        </p>
        <div class="card-wrapper">
            {#if cards.length > 1}
                <button class="scroll" onclick={scrollLeft}
                    ><img src={ScrollLeft} alt="Scroll left" /></button
                >
            {/if}
            {#if cards.length >= 1}
                <div class="card">
                    <button
                        class="card-btn"
                        onclick={() => selectCard(card.id)}
                    >
                        <!-- this inline styling fixes Svelte rendering two texts while outro is playing -->
                        <div
                            class="card-content"
                            style="position: relative; width: 100%; height: 100%;"
                        >
                            {#key card}
                                <p
                                    id="title"
                                    in:fade|local={{
                                        duration: 200,
                                        delay: 200,
                                    }}
                                    out:fade|local={{ duration: 200 }}
                                    style="position: absolute; transform: translateY(-120%);"
                                >
                                    {card.title}
                                </p>
                                <p
                                    id="subject"
                                    in:fade|local={{
                                        duration: 200,
                                        delay: 200,
                                    }}
                                    out:fade|local={{ duration: 200 }}
                                    style="position: absolute; transform: translateY(120%);"
                                >
                                    {card.subject}
                                </p>
                            {/key}
                        </div></button
                    >
                </div>
            {/if}
            {#if cards.length === 0}
                <p style="margin-bottom: 40px;">{$_("home.no_cards")}</p>
            {/if}
            {#if cards.length > 1}
                <button class="scroll" onclick={scrollRight}
                    ><img src={ScrollRight} alt="Scroll left" /></button
                >
            {/if}
        </div>
        {#if cards.length >= 1}
            <div id="buttons">
                <button
                    id="delete"
                    onclick={() => {
                        deleteCard(card.id);
                    }}
                    class="icon-btn"
                >
                    <img src={Delete} alt="" style="width: 20px" />
                </button>
                <button id="edit" onclick={edit} class="icon-btn">
                    <img src={Edit} alt="" style="width: 20px" />
                </button>
            </div>
        {/if}
        <button class="new" onclick={newFunc}>{$_("home.new_card")}</button>
    </div>
</div>

<style>
    @font-face {
        font-family: Inter;
        src: url("../fonts/Inter.ttf");
    }

    * {
        color: white;
        font-family: Inter;
    }

    *:focus {
        outline: 2px solid #535353;
    }

    button {
        border: none;
        border-radius: 50px;
        background-color: #373737;
        padding: 20px;
        font-size: 13pt;
        cursor: pointer;
    }

    button:disabled {
        opacity: 50%;
        cursor: not-allowed;
    }

    button:hover:not(:disabled, .card-btn) {
        background-color: #535353;
    }

    .wrapper {
        width: 100%;
        margin: 0;
        padding: 0;
    }

    .main {
        width: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    .text {
        font-size: 20pt;
    }

    .card-wrapper {
        display: flex;
        justify-content: center;
        align-items: center;
    }

    .card {
        min-width: 375px;
        min-height: 250px;
        max-width: 375px;
        max-height: 250px;
        background-color: #373737;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        border-radius: 20px;
        margin: 30px 30px;
    }

    #buttons {
        display: flex;
    }

    .icon-btn {
        height: 40px;
        width: 40px;
        border: #535353 3px solid;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    #delete {
        margin-right: 10px;
    }

    #edit {
        margin-left: 10px;
    }

    .card-btn {
        width: 100%;
        height: 225px;
    }:focus {
        outline: none;
    }

    .card-content {
        width: 100%;
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }

    @keyframes slide-left {
        100% {
            transform: translateX(-100%);
            opacity: 0;
        }
    }
    @keyframes slide-right {
        0% {
            transform: translateX(100%);
            opacity: 0;
        }

        100% {
            transform: translateX(0%);
            opacity: 100%;
        }
    }

    #title {
        font-size: x-large;
        font-weight: 590;
    }

    #subject {
        font-weight: 350;
    }

    .scroll {
        height: 70px;
        width: 70px;
        border: #535353 3px solid;
    }

    .new {
        margin-top: 60px;
    }
</style>
