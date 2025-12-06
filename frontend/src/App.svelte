<script>
  // @ts-ignore
  import logo from "./assets/images/logo.png";
  // @ts-ignore
  import HomeImg from "./assets/images/home.png";
  import SettingsImg from "./assets/images/settings.png";

  import Home from "./assets/components/Home.svelte";
  import New from "./assets/components/New.svelte";
  import Card from "./assets/components/Card.svelte";
  import End from "./assets/components/End.svelte";

  import { register, init, isLoading, _, locale, locales } from "svelte-i18n";
  import { onMount } from "svelte";

  // @ts-ignore
  import { ReadSettings, SaveSettings } from "../wailsjs/go/main/App";
    import { fly } from "svelte/transition";

  register("en", () => import("./assets/locales/en.json"));
  register("it", () => import("./assets/locales/it.json"));

  init({
    fallbackLocale: "en",
    initialLocale: "en",
  });

  let selectedLang = $state("");

  onMount(async () => {
    let settings = await ReadSettings();
    console.log("Settings:", settings);
    selectedLang = settings.language;
    init({
      fallbackLocale: "en",
      initialLocale: settings.language,
    });
    console.log("Translations applied! Can now render.");
  });

  locale.subscribe(async (lang) => {
    let settings = await ReadSettings();
    settings.language = lang;
    await SaveSettings(settings);
  });

  let current = $state("home");
  let cardId = $state("");

  let rerender = $state(0);

  let showingSettings = $state(false);

  function newFn() {
    current = "new";
  }

  function back() {
    current = "home";
  }

  function card(id) {
    cardId = id;
    current = "card";
  }

  function finish() {
    current = "end";
  }

  function correct() {
    current = "correction";
  }

  function del() {
    rerender += 1;
  }

  function edit(id) {
    cardId = id;
    current = "edit";
  }
</script>

{#if $isLoading}
  <h1>Loading...</h1>
{:else}
  <main>
    <button onclick={back} id="logo-btn"
      ><img alt="Wails logo" id="logo" src={logo} /></button
    >
    <button id="home-btn" onclick={back}>
      <img src={HomeImg} alt="Home" />
    </button>
    <button
      id="settings-btn"
      onclick={() => {
        showingSettings = !showingSettings;
      }}
    >
      <img src={SettingsImg} alt="Home" />
    </button>
    {#if showingSettings}
      <div id="settings" transition:fly={{ duration: 250, y: -20 }}>
        <div id="lang-div">
          <label for="lang">{$_("language")}</label>
          <select
            name="lang"
            id="lang"
            bind:value={selectedLang}
            onchange={() => {
              locale.set(selectedLang);
            }}
          >
            {#each $locales as loc}
              <option value={loc}>{loc}</option>
            {/each}
          </select>
        </div>
      </div>
    {/if}
    <div id="content">
      {#if current === "home"}
        {#key rerender}
          <Home onNew={newFn} onSelect={card} onDelete={del} onEdit={edit}
          ></Home>
        {/key}
      {:else if current === "new"}
        <New onDone={back}></New>
      {:else if current === "card"}
        <Card id={cardId} onFinish={finish}></Card>
      {:else if current === "end"}
        <End id={cardId} onStart={correct}></End>
      {:else if current === "correction"}
        <Card id={cardId} withCorrection={true} onFinish={finish}></Card>
      {:else if current === "edit"}
        <New onDone={back} edit={true} id={cardId}></New>
      {/if}
    </div>
  </main>
{/if}

<style>
  main {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    margin-bottom: 20px;
  }

  select {
    appearance: none;
    background-color: #323232;
    border: 2px solid #535353;
    color: white;
    width: 50px;
    text-align: center;
  }

  option {
    text-align: center;
  }

  select:focus {
    outline: none;
  }

  #logo {
    width: 30vw;
    margin-top: 2vh;
  }

  #logo-btn {
    border: none;
    background-color: transparent;
    cursor: pointer;
  }

  #home-btn {
    border-radius: 50px;
    background-color: #373737;
    padding: 20px;
    font-size: 13pt;
    cursor: pointer;
    height: 40px;
    width: 40px;
    border: #535353 3px solid;
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    top: 15px;
    left: 15px;
  }

  #home-btn:hover {
    background-color: #535353;
  }

  #settings-btn {
    border-radius: 50px;
    background-color: #373737;
    padding: 20px;
    font-size: 13pt;
    cursor: pointer;
    height: 40px;
    width: 40px;
    border: #535353 3px solid;
    display: flex;
    align-items: center;
    justify-content: center;
    position: absolute;
    top: 15px;
    right: 15px;
  }

  #settings-btn:hover {
    background-color: #535353;
  }

  #settings {
    background-color: #323232;
    border-radius: 20px;
    border: #535353 3px solid;
    display: flex;
    flex-direction: column;
    position: absolute;
    right: 0;
    top: 50px;
    width: 300px;
    margin: 20px;
    padding: 15px;
  }

  #lang-div {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  #content {
    border: 2px solid #323232;
    border-radius: 20px;
    flex: 1;
    width: calc(100vw - 40px);
    margin: 20px;
  }
</style>
