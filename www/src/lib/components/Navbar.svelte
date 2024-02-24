<script>
    import logo from '$lib/assets/logo.png';
</script>

<div class="navbar bg-base-200 rounded-full my-5">
    <div class="navbar-start">
        <div class="btn btn-ghost btn-circle avatar">
            <div class="w-10 rounded-full">
                <img alt="Motmot" src={logo} />
            </div>
        </div>

        <p class="ml-1 cursor-default text-2xl">Motmot</p>
    </div>

    <div class="navbar-end lg:flex">
        <ul class="menu menu-horizontal px-1 mr-1">
            <li><a href="/groups" class="rounded-full">Courses</a></li>
        </ul>
    </div>
</div>