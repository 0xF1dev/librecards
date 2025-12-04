<script>
    import { onMount } from "svelte";

    // @ts-ignore
    import { GetErrors, GetCard } from "../../../wailsjs/go/main/App";
    // @ts-ignore
    import ScrollLeft from "../images/scroll-left.png";
    // @ts-ignore
    import ScrollRight from "../images/scroll-right.png";

    import { fly } from "svelte/transition";
    import SvelteMarkdown from "@humanspeak/svelte-markdown";

    let { id, onStart } = $props();

    let errors;
    let title = $state("");
    let card;

    let section = 0;
    let question = 0;

    let sectionTitle = $state("");
    let text = $state("");

    let inAnimation = $state({ duration: 250, x: 40, delay: 250 });
    let outAnimation = $state({ duration: 250, x: -40 });

    let hasErrors = $state(false);

    onMount(async () => {
        errors = await GetErrors(id);
        console.log(errors);

        if (errors.errors.length == 0) {
            return;
        } else {
            hasErrors = true;
        }

        title = errors.title;

        let fullCard = await GetCard(id);

        card = structuredClone(fullCard);
        card.data.sections = [];

        console.log("CARD: ", fullCard);
        console.log("ERRORS: ", errors);

        errors.errors.forEach((element) => {
            let section = {
                title: fullCard.data.sections[element.section].title,
                questions: [],
            };
            element.questions.forEach((question) => {
                section.questions.push(
                    fullCard.data.sections[element.section].questions[question],
                );
            });
            card.data.sections = [...card.data.sections, section];
        });

        console.log(card);

        sectionTitle = card.data.sections[0].title;
        text = card.data.sections[0].questions[0].question;
    });

    function scrollRight() {
        console.log(card);

        inAnimation = { duration: 250, x: 40, delay: 250 };
        outAnimation = { duration: 250, x: -40 };

        if (question >= card.data.sections[section].questions.length - 1) {
            if (section >= card.data.sections.length - 1) {
                section = 0;
            } else {
                section += 1;
            }
            question = 0;
        } else {
            question += 1;
        }

        sectionTitle = card.data.sections[section].title;
        text = card.data.sections[section].questions[question].question;
    }

    function scrollLeft() {
        console.log(card);

        inAnimation = { duration: 250, x: -40, delay: 250 };
        outAnimation = { duration: 250, x: 40 };

        if (question == 0) {
            if (section <= 0) {
                return;
            }
            section -= 1;
            question = card.data.sections[section].questions.length - 1;
        } else {
            question -= 1;
        }

        sectionTitle = card.data.sections[section].title;
        text = card.data.sections[section].questions[question].question;
    }

    function start() {
        onStart();
    }
</script>

<div class="wrapper">
    {#if hasErrors}
        <h1 bind:innerText={title} contenteditable="false" id="title">
            Loading...
        </h1>
        <button id="correct" onclick={start}>Start correction</button>
        <div id="preview-wrapper">
            <div class="section-title" style="display: grid;">
                {#key sectionTitle}
                    <p
                        id="section-title"
                        bind:innerText={sectionTitle}
                        contenteditable="false"
                        in:fly|local={{ duration: 250, y: 20, delay: 250 }}
                        out:fly|local={{ duration: 250, y: -20 }}
                        style="grid-area: 1/1;"
                    >
                        Loading...
                    </p>
                {/key}
            </div>
            <div class="card-wrapper">
                <div class="card-interactions">
                    <button class="scroll" id="previous" onclick={scrollLeft}
                        ><img src={ScrollLeft} alt="Go back" /></button
                    >
                    <div class="card" id="card">
                        <!-- position: absolute; fixes Svelte rendering two texts while outro is playing -->
                        {#key text}
                            <div
                                in:fly|local={inAnimation}
                                out:fly|local={outAnimation}
                                style="position: absolute;"
                            >
                                <SvelteMarkdown source={text} />
                            </div>
                        {/key}
                    </div>
                    <button class="scroll" id="next" onclick={scrollRight}
                        ><img src={ScrollRight} alt="Go forward" /></button
                    >
                </div>
            </div>
        </div>
    {:else}
        <h1 id="no-errs">No errors!</h1>
    {/if}
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

    button:hover:not(:disabled) {
        background-color: #535353;
    }

    .wrapper {
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    #no-errs {
        margin-top: 50px;
        font-weight: 500;
    }

    #preview-wrapper {
        margin-top: 20px;
    }

    #section-title {
        font-size: xx-large;
    }

    .card-wrapper {
        flex: 1 1 auto;
    }

    .card-interactions {
        display: flex;
        align-items: center;
    }

    .scroll {
        height: 70px;
        width: 70px;
        border: #535353 3px solid;
    }

    .card {
        min-width: 400px;
        min-height: 280px;
        background-color: #373737;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        border-radius: 20px;
        margin: 0 20px;
    }
</style>
