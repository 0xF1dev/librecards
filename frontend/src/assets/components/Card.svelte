<script>
    // @ts-nocheck
    import ScrollLeft from "../images/scroll-left.png";
    import ScrollRight from "../images/scroll-right.png";
    import End from "../images/end.png";

    // @ts-ignore
    import {
        GetCard,
        SaveError,
        SaveCorrect,
        GetErrors,
    } from "../../../wailsjs/go/main/App";
    import { onMount } from "svelte";
    import { fly } from "svelte/transition";

    import SvelteMarkdown from "@humanspeak/svelte-markdown";

    import { _ } from "svelte-i18n";

    let { id, withCorrection = false, onFinish } = $props();

    let card;

    let section = 0;
    let question = 0;

    let realValues = [];

    let sectionTitle = $state("");
    let text = $state("");

    let cardTitle = $state("");

    let showingAnswer = $state(false);

    let correct = $state(false);
    let wrong = $state(false);
    let lastStatus = "";

    let inAnimation = $state({ duration: 250, x: 40, delay: 250 });
    let outAnimation = $state({ duration: 250, x: -40 });

    let scrollImage = $state(ScrollRight);

    let canScrollLeft = $state(false);

    onMount(async () => {
        if (withCorrection == true) {
            let errors = await GetErrors(id);
            console.log(errors);

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
                        fullCard.data.sections[element.section].questions[
                            question
                        ],
                    );
                });
                card.data.sections = [...card.data.sections, section];
                realValues.push(element);
            });

            console.info("card: " + card);
            console.info("real values: " + realValues)

            sectionTitle = card.data.sections[0].title;
            text = card.data.sections[0].questions[0].question;
        } else {
            card = await GetCard(id);
            console.log(card);

            sectionTitle = card.data.sections[section].title;
            text = card.data.sections[section].questions[question].question;
        }

        cardTitle = card.data.title + " (#" + card.id + ")";

        if (
            question >= card.data.sections[section].questions.length - 1 &&
            section >= card.data.sections.length - 1
        ) {
            scrollImage = End;
        } else {
            scrollImage = ScrollRight;
        }
    });

    function showAnswer() {
        if (showingAnswer === false) {
            inAnimation = { duration: 250, y: 20, delay: 250 };
            outAnimation = { duration: 250, y: -20 };

            text = card.data.sections[section].questions[question].answer;
            showingAnswer = true;
            if (lastStatus !== "") {
                lastStatus === "correct" ? (correct = true) : (wrong = true);
            }
        } else {
            inAnimation = { duration: 250, y: 20, delay: 250 };
            outAnimation = { duration: 250, y: -20 };

            text = card.data.sections[section].questions[question].question;
            showingAnswer = false;
            correct = false;
            wrong = false;
        }
    }

    async function onCorrect() {
        console.log("saving correction");
        if (withCorrection == true) {
            let err = await SaveCorrect(id, realValues[section].section, realValues[section].questions[question]);
            if (err === 1) {
                console.log("error!");
                return;
            }
        } else {
            let err = await SaveCorrect(id, section, question);
            if (err === 1) {
                console.log("error!");
                return;
            }
        }

        wrong = false;
        correct = true;
        lastStatus = "correct";
    }

    async function onWrong() {
        console.log("saving err");
        if (withCorrection == true) {
            let err = await SaveError(id, realValues[section].section, realValues[section].questions[question]);
            if (err === 1) {
                console.log("error!");
                return;
            }
        } else {
            let err = await SaveError(id, section, question);
            if (err === 1) {
                console.log("error!");
                return;
            }
        }

        correct = false;
        wrong = true;
        lastStatus = "wrong";
    }

    function scrollRight() {
        console.log(card);

        inAnimation = { duration: 250, x: 40, delay: 250 };
        outAnimation = { duration: 250, x: -40 };

        correct = false;
        wrong = false;
        lastStatus = "";

        showingAnswer = false;

        if (question >= card.data.sections[section].questions.length - 1) {
            if (section >= card.data.sections.length - 1) {
                onFinish();
                return;
            }
            question = 0;
            section += 1;
        } else {
            question += 1;
        }

        sectionTitle = card.data.sections[section].title;
        text = card.data.sections[section].questions[question].question;
        canScrollLeft = true;

        if (
            question >= card.data.sections[section].questions.length - 1 &&
            section >= card.data.sections.length - 1
        ) {
            scrollImage = End;
        } else {
            scrollImage = ScrollRight;
        }
    }

    function scrollLeft() {
        console.log(card);

        inAnimation = { duration: 250, x: -40, delay: 250 };
        outAnimation = { duration: 250, x: 40 };

        correct = false;
        wrong = false;
        lastStatus = "";

        showingAnswer = false;

        if (question == 0) {
            if (section <= 0) {
                return;
            }
            section -= 1;
            question = card.data.sections[section].questions.length - 1;
        } else {
            question -= 1;
            if (section == 0 && question == 0) {
                canScrollLeft = false;
            }
        }

        sectionTitle = card.data.sections[section].title;
        text = card.data.sections[section].questions[question].question;

        if (
            question >= card.data.sections[section].questions.length - 1 &&
            section >= card.data.sections.length - 1
        ) {
            scrollImage = End;
        } else {
            scrollImage = ScrollRight;
        }
    }
</script>

<div class="wrapper" id="wrapper">
    <div class="main">
        <!-- grid display fixes Svelte rendering two texts while outro is playing -->
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
                <button
                    class="scroll"
                    id="previous"
                    disabled={!canScrollLeft}
                    onclick={scrollLeft}
                    ><img src={ScrollLeft} alt="Go back" /></button
                >
                <div class="card" id="card" class:correct class:wrong>
                    <!-- position: absolute; fixes Svelte rendering two texts while outro is playing -->
                    {#key text}
                        <div
                            in:fly|local={inAnimation}
                            out:fly|local={outAnimation}
                            onoutroend={() => {
                                inAnimation = {
                                    duration: 250,
                                    x: 40,
                                    delay: 250,
                                };
                                outAnimation = { duration: 250, x: -40 };
                            }}
                            style="position: absolute;"
                        >
                            <SvelteMarkdown source={text} />
                        </div>
                    {/key}
                </div>
                <button class="scroll" id="next" onclick={scrollRight}
                    ><img src={scrollImage} alt="Go forward" /></button
                >
            </div>
            <div class="buttons">
                {#if showingAnswer === true}
                    <button id="wrong" onclick={onWrong}>{$_('card.wrong')}</button>
                {/if}
                <button id="reveal" onclick={showAnswer}
                    >{#if showingAnswer === true}
                        {$_('card.hide_answer')}
                    {:else}
                        {$_('card.show_answer')}
                    {/if}</button
                >
                {#if showingAnswer === true}
                    <button id="correct" onclick={onCorrect}>{$_('card.correct')}</button>
                {/if}
            </div>
        </div>
        <div class="card-title">
            <p
                id="card-title"
                bind:innerText={cardTitle}
                contenteditable="false"
            >
                Loading...
            </p>
        </div>
    </div>
</div>
<div class="error" id="error">
    <p class="error-msg">{$_('card.error')}</p>
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

    .main {
        height: 100%;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
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
        min-width: 500px;
        min-height: 350px;
        background-color: #373737;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        border-radius: 20px;
        margin: 30px 30px 60px;
    }

    .buttons {
        display: flex;
        justify-content: center;
        align-items: center;
    }

    #reveal {
        margin-left: 20px;
        margin-right: 20px;
    }

    #wrong {
        background-color: #5a3434;
    }

    #wrong:hover {
        background-color: #814d4d;
    }

    #correct {
        background-color: #345a36;
    }

    #correct:hover {
        background-color: #4b7a4e;
    }

    .card-title {
        margin: 20px 0 5px;
    }

    #card-title {
        color: #373737;
    }

    .error {
        display: none;
        align-items: center;
        justify-content: center;
    }

    .error-msg {
        font-size: 30pt;
    }

    @keyframes wrong {
        100% {
            background-color: #5a3434;
        }
    }

    @keyframes correct {
        100% {
            background-color: #345a36;
        }
    }

    .card.wrong {
        animation: wrong 0.5s forwards;
    }

    .card.correct {
        animation: correct 0.5s forwards;
    }
</style>
