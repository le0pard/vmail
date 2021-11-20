<div class="empty-report">
  <div class="empty-report-body">
    <h3>VMail - check HTML & CSS compatibility with email clients</h3>
    <p class="empty-report-text">
      Email clients use different rendering standards. This is why your email
      can be displayed not as you designed it. You need to check that your
      message code won't cause rendering issues.
    </p>
    <p class="empty-report-text">
      Vmail (Validate Email) check HTML & CSS email template content in search
      of problematic elements. For each it finds, it displays the list of email
      clients that lack support for it or support it only partially.
    </p>
    <p class="empty-report-text">
      Submit your HTML email template and get report about it is possible issues
      in different email clients
    </p>
  </div>
</div>

<style>
  .empty-report {
    border: 0;
    box-shadow: 0 10px 35px 0 rgb(56 71 109 / 8%);
    position: relative;
    display: flex;
    flex-direction: column;
    min-width: 0;
    background-color: var(--cardBgColor);
    border-radius: 0.4rem;
  }

  .empty-report-body {
    flex: 1 1 auto;
    padding: 2rem;
  }

  .empty-report-text {
    color: var(--cardBaseColor);
  }
</style>
