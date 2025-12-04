<script>
  // @ts-ignore
  import logo from "./assets/images/logo.png";
  // @ts-ignore
  import HomeImg from "./assets/images/home.png";

  import Home from "./assets/components/Home.svelte";
  import New from "./assets/components/New.svelte";
  import Card from "./assets/components/Card.svelte";
  import End from "./assets/components/End.svelte";

  let current = $state("home");
  let cardId = $state("");

  let rerender = $state(0);

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
    cardId = id
    current = "edit"
  }
</script>

<main>
  <button onclick={back} id="logo-btn"
    ><img alt="Wails logo" id="logo" src={logo} /></button
  >
  <button id="home-btn" onclick={back}>
    <img src={HomeImg} alt="Home" />
  </button>
  <div id="content">
    {#if current === "home"}
      {#key rerender}
        <Home onNew={newFn} onSelect={card} onDelete={del} onEdit={edit}></Home>
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

<style>
  main {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    /* padding-bottom: 20px; */
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

  #content {
    border: 2px solid #323232;
    border-radius: 20px;
    flex: 1;
    width: calc(100vw - 40px);
    margin: 20px;
  }
</style>
