<script>
    // @ts-ignore
    import Add from "../images/add.png";
    // @ts-ignore
    import Delete from "../images/delete.png";
    // @ts-ignore
    import Done from "../images/done.png";

    // @ts-ignore
    import { NewCard, GetCard, UpdateCard } from "../../../wailsjs/go/main/App";
    import { onMount } from "svelte";

    let { onDone, edit = false, id = "" } = $props();

    let title = $state("");
    let subject = $state("");
    let sections = $state([
        { title: "", questions: [{ question: "", answer: "" }] },
    ]);

    let error = $state(0);

    onMount(async () => {
        if (edit === true) {
            let card = await GetCard(id);
            title = card.data.title;
            subject = card.data.subject;
            sections = card.data.sections;
        }
    });

    function addSection() {
        sections = [
            ...sections,
            { title: "", questions: [{ question: "", answer: "" }] },
        ];
    }

    function removeSection(index) {
        if (sections.length - 1 !== 0) {
            sections = sections.filter((_, i) => i !== index);
        }
    }

    function addQuestion(sectionIndex) {
        sections[sectionIndex].questions = [
            ...sections[sectionIndex].questions,
            { question: "", answer: "" },
        ];
    }

    function removeQuestion(sectionIndex, index) {
        if (sections[sectionIndex].questions.length - 1 !== 0) {
            sections[sectionIndex].questions = sections[
                sectionIndex
            ].questions.filter((_, i) => i !== index);
        }
    }

    async function createCard() {
        let newCard = { title: title, subject: subject, sections: sections };
        console.log(newCard);
        if (edit === true) {
            // @ts-ignore
            let error = await UpdateCard(id, newCard);
            if (error === 0) {
                onDone();
            }
            return
        }

        // @ts-ignore
        error = await NewCard(newCard);
        if (error === 0) {
            onDone();
        }
    }
</script>

<div class="wrapper">
    <input
        type="text"
        name=""
        id="title"
        placeholder="Title"
        bind:value={title}
    />
    <input
        type="text"
        name=""
        id="subject"
        placeholder="Subject"
        bind:value={subject}
    />
    <div class="sections">
        {#each sections as section, i}
            <div class="section">
                <input
                    type="text"
                    name=""
                    id="section-title"
                    placeholder="Section title ({i + 1})"
                    bind:value={section.title}
                />
                <div class="questions">
                    <ul>
                        {#each sections[i].questions as question, iQuestion}
                            <li>
                                <div class="question-wrapper">
                                    <textarea
                                        name=""
                                        id="question"
                                        placeholder="Question"
                                        bind:value={question.question}
                                    ></textarea>
                                    {#if sections[i].questions.length - 1 !== 0}
                                        <button
                                            class="delete icon-btn"
                                            id="delete"
                                            onclick={() => {
                                                removeQuestion(i, iQuestion);
                                            }}
                                            ><img
                                                src={Delete}
                                                alt="Delete question"
                                            /></button
                                        >
                                    {:else}
                                        <span style="width: 45px;"></span>
                                    {/if}
                                    <textarea
                                        name=""
                                        id="answer"
                                        placeholder="Answer"
                                        bind:value={question.answer}
                                    ></textarea>
                                </div>
                            </li>
                        {/each}
                    </ul>
                </div>
                <div class="btns">
                    <button
                        class="add icon-btn"
                        id="add"
                        onclick={() => {
                            addQuestion(i);
                        }}><img src={Add} alt="New question" /></button
                    >
                    {#if sections.length - 1 !== 0}
                        <button
                            class="delete-section icon-btn"
                            id="delete"
                            onclick={() => {
                                removeSection(i);
                            }}><img src={Delete} alt="Delete section" /></button
                        >
                    {/if}
                </div>
            </div>
        {/each}
    </div>
    <button class="add-section" onclick={addSection}>Add section</button>
    <button class="done icon-btn" id="done" onclick={createCard}
        ><img src={Done} alt="Done" /></button
    >
    {#if error === 1}
        <h1>Error!</h1>
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

    input[type="text"],
    textarea {
        border: none;
        border-radius: 50px;
        background-color: #373737;
        resize: none;
    }

    textarea {
        padding-left: 20px;
        padding-right: 20px;
        padding-top: 10px;
        border-radius: 30px;
        margin: 0 10px;
    }

    ul {
        list-style: none;
        padding: 0;
    }

    .wrapper {
        display: flex;
        flex-direction: column;
        align-items: center;
    }

    #title {
        font-size: 30pt;
        text-align: center;
        margin: 20px;
    }

    #subject {
        font-size: 15pt;
        text-align: center;
        margin-bottom: 40px;
    }

    .section {
        display: flex;
        flex-direction: column;
        align-items: center;
        border-top: #535353 2px solid;
        border-bottom: #535353 2px solid;
        padding: 15px 0 15px;
        width: 40vw;
    }

    #section-title {
        font-size: 15pt;
        padding-left: 10px;
    }

    .icon-btn {
        height: 40px;
        width: 40px;
        border: #535353 3px solid;
        display: flex;
        align-items: center;
        justify-content: center;
    }

    .done {
        height: 70px;
        width: 70px;
    }

    .question-wrapper {
        display: flex;
        align-items: center;
    }

    .question-wrapper:not(:nth-child(0)) {
        margin-top: 10px;
    }

    .btns {
        display: flex;
    }

    .delete-section {
        margin-left: 20px;
    }

    .add-section {
        margin: 40px;
    }
</style>
