# Graph Report - app-listener  (2026-11-07)

## Corpus Check
- 575 files · ~1,821,003 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7550 nodes · 34969 edges · 267 communities (215 shown, 52 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 3102 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `cff5c42b`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- render.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- GuardEvent
- ConfinedEntry
- RefreshSection
- Guard
- buildOneGuard
- DaemonEvent
- IntegrationSuite
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- DaemonUseCase
- auth.go
- op
- guard.bpf.c
- Finder
- inode_dev
- gui/session.go
- rcUser
- editView
- BinaryEntry
- .apply
- sync.Mutex
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- systemd.go
- NetEventType
- root.go
- UI
- fileedit_model_test.go
- monitorUnitTest
- pinstate.go
- updatePage
- github.com/charmbracelet/bubbletea.Cmd
- stdio
- runView
- fadeScroll
- time.Time
- protectedPage
- Client
- highlight_test.go
- fyne.io/fyne/v2/driver/desktop.Cursor
- applier
- backend_test.go
- Serve
- GuardInodeKey
- IntegrationSuite
- installPage
- deploy
- Ledger
- golang.org/x/sys/unix.Stat_t
- daemon.conf reference
- unix_connect.c
- IntegrationSuite
- wire.go
- trustManager
- parseDaemonEvents
- Vault
- CandidateDir
- changelog.go
- fcntl
- FileTarget
- FileEvent
- codecheck/main.go
- NetGuard
- liveSession
- process_vm_readv.c
- .filterEvent
- safeio.go
- watchSet
- daemon_multicall_test.go
- fakeGuardRepo
- bpfstats/main.go
- exe_supersede.h
- loginServer
- compileGlobNames
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LaunchRule
- User
- netinject.c
- dialLiveSession
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- recordFlow
- journal
- Well-known bypass classes
- IntegrationSuite
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- vaultfs.go
- Build & sign release assets (reusable workflow)
- planUpdaters
- newTestTarget
- Gate
- UnifiedDiff
- NetworkMonitor
- IntegrationSuite
- reloadOnce
- Row
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- catalogPatterns
- IntegrationSuite
- WebSocket /ws client
- sync.Once
- inet_connect.c
- install.sh
- validateConfigText
- WithSuspectsFrom
- make test-integration (rootful Docker suite)
- SelectOrphans
- .exportTable
- github.com/spf13/cobra.Command
- trace-app-libs.sh
- IntegrationSuite
- guardModel
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- session
- KernelDev
- Monitor
- Resource
- TestMonitorUseCaseLifecycle
- maintainerKeeper
- btrfsLayoutFrom
- Recorder
- SampleContent
- fakeStore
- Backup
- IntegrationSuite
- .ExpandLibDirWriters
- go_pkg_strings
- copytree.go
- verdict
- Installation
- NetGuardEvent
- graphify-refresh.sh
- .Layout
- backend.go
- TestCleanRel
- serverAuth
- feed.go
- shell
- classify
- sources_test.go
- os.FileInfo
- searchDirs
- sources
- Limitations
- hashCache
- plan_test.go
- net_tester/main.go
- editauth_test.go
- mustSubkey
- tempGrantSetup
- IntegrationSuite
- New
- TrustGuard
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- IntegrationSuite
- os.File
- Self-updating apps and the catalog
- run.go
- WithBottomBar
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- IntegrationSuite
- pathCache
- ensurePortAnchors
- auditGroupCoverage
- .netIntegrityContainer
- changelog_test.go
- SystemTrusted
- monitor
- IntegrationSuite
- absPath
- .admitGated
- runMonitor
- catalogRefresher
- update_test.go
- .populateWatchInodes
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- runNetworkGuard
- runGenKey
- find_test.go
- appTheme
- ResolvePinBase
- IntegrationSuite
- inodeChain
- batchInodeRows
- supersede_probe.c
- decryptStep
- bintrust_test.go
- launch_probe.c
- IntegrationSuite
- github.com/Virgula0/app-listener/tools/codecheck
- procstats.go
- taint_victim.c
- secretJudge
- mutedText
- github.com/charmbracelet/bubbles/textarea.Model
- Troubleshooting
- mountTable
- GenerateConf
- .runEditorHarness
- IntegrationSuite
- testLimiter
- linkView
- bpfLsmListed
- downloadFile
- Development
- app-listener documentation
- ShortRefused
- run
- VouchedDevs
- launchKey
- Vault
- CleanupStalePins
- recorder/plan.go
- monitor.md
- network monitor: see what a program talks to
- .eventTypeAllowed
- daemon-config.md
- noteReloadedSettings
- SecretMatch
- `[launch_rules]`
- harnessSuite
- bleedLayout
- parsePasswd
- LiveEmptyWhitelistRejected
- daemonconfig.go

## God Nodes (most connected - your core abstractions)
1. `Config` - 149 edges
2. `Guard` - 126 edges
3. `GuardInodeKey` - 116 edges
4. `Load()` - 101 edges
5. `absPath()` - 97 edges
6. `shell` - 97 edges
7. `IntegrationSuite` - 83 edges
8. `writeConfig()` - 79 edges
9. `fileEditModel` - 77 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Modes` --references--> `monitor()`  [INFERRED]
  docs/README.md → integrationtests/secret_probe/main.go
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `Workspaces` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `eBPF at three levels` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (267 total, 52 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.15
Nodes (32): tempJournalRow, argvMatch, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion (+24 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (189): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+181 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (7): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig()

### Community 3 - "render.go"
Cohesion: 0.12
Nodes (13): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+5 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (85): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+77 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.09
Nodes (51): prctl_gate(), guard_task_prctl(), admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader() (+43 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.21
Nodes (6): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent()

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (22): addErrKind, ReplacementCheck, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+14 more)

### Community 9 - "mkdirs"
Cohesion: 0.12
Nodes (38): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetHomeDepartureTriggersRefresh(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds() (+30 more)

### Community 10 - "GuardEvent"
Cohesion: 0.07
Nodes (18): holdHeadless(), HeadlessLine(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary() (+10 more)

### Community 11 - "ConfinedEntry"
Cohesion: 0.13
Nodes (29): binaryVerifyState, HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats() (+21 more)

### Community 12 - "RefreshSection"
Cohesion: 0.10
Nodes (28): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+20 more)

### Community 13 - "Guard"
Cohesion: 0.04
Nodes (27): deferredBinary, exeRow, keptEntries, rootHandle, supersededBinary, tempMaskOp, tempRow, openRoot() (+19 more)

### Community 14 - "buildOneGuard"
Cohesion: 0.09
Nodes (35): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), selfProtectSpecs(), startSelfGuard(), displayEntries(), resolveGuardConfig(), runGuard() (+27 more)

### Community 15 - "DaemonEvent"
Cohesion: 0.10
Nodes (28): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless() (+20 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.08
Nodes (69): argv_env(), argv_matched(), argv_not_root(), argv_rule_at(), argv_short(), argv_word(), caller_updates(), creds_overridden() (+61 more)

### Community 18 - "navItem"
Cohesion: 0.03
Nodes (25): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panel, panelRenderer (+17 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.13
Nodes (26): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+18 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.14
Nodes (50): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGrantEditAccessWidensNarrows() (+42 more)

### Community 22 - "DaemonUseCase"
Cohesion: 0.07
Nodes (12): GuardRepository, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots(), TestPartitionEncryptionRootsSplitsDirsAndFiles(), TestRepeatLogOncePerErrorAndOnRecovery(), uniqueEncryptionRoots() (+4 more)

### Community 23 - "auth.go"
Cohesion: 0.13
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "op"
Cohesion: 0.16
Nodes (13): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+5 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.07
Nodes (97): exe_refused(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), bus_key(), check_and_emit_args(), chmod_only_drops_write(), code_suspect() (+89 more)

### Community 26 - "Finder"
Cohesion: 0.16
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 27 - "inode_dev"
Cohesion: 0.43
Nodes (4): ctx_ptr(), inode_dev(), sb_dev(), guard_inode_free()

### Community 28 - "gui/session.go"
Cohesion: 0.07
Nodes (40): invokingIDs(), ancestorEnv(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+32 more)

### Community 29 - "rcUser"
Cohesion: 0.12
Nodes (42): rcTarget, TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir() (+34 more)

### Community 30 - "editView"
Cohesion: 0.10
Nodes (9): codeEntry, editView, handButton, Entry, fileSize(), joinRel(), newCodeEntry(), relOf() (+1 more)

### Community 31 - "BinaryEntry"
Cohesion: 0.11
Nodes (12): BinaryEntry, diffBinaries(), TestDiffBinaries(), netGuardPill(), NetworkGuardRepository, clonePolicy(), NetEngines, NetExtra (+4 more)

### Community 32 - ".apply"
Cohesion: 0.11
Nodes (20): ValidateTargets(), fixtureExe(), mustParse(), privateTempRoots(), record(), TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything() (+12 more)

### Community 33 - "sync.Mutex"
Cohesion: 0.10
Nodes (23): readConfigPut(), TestReadGrantRequestConfig(), authSelfPeer(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+15 more)

### Community 34 - "picker"
Cohesion: 0.07
Nodes (32): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, lstatOf(), outsideHardlinks() (+24 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.18
Nodes (31): busPath(), connexec(), connFD(), createItem(), dial(), dupkey(), execInto(), failf() (+23 more)

### Community 37 - "Message"
Cohesion: 0.10
Nodes (15): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+7 more)

### Community 38 - "systemd.go"
Cohesion: 0.07
Nodes (52): offerDaemonStart(), applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall(), installBinaryOnly(), markInstallComplete(), runRefreshCatalogOnly(), softenAutomatedRefreshErr() (+44 more)

### Community 39 - "NetEventType"
Cohesion: 0.08
Nodes (15): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), newNetEventChecks(), FormatAddr(), NetEventTypes() (+7 more)

### Community 40 - "root.go"
Cohesion: 0.09
Nodes (14): installHooks(), releaseBuild(), Run(), staged(), closeDumpFile(), postRun(), Progress, User (+6 more)

### Community 41 - "UI"
Cohesion: 0.08
Nodes (23): showChangelog(), channelOf(), confirmUpdate(), isTerminal(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses() (+15 more)

### Community 42 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 44 - "pinstate.go"
Cohesion: 0.11
Nodes (32): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder() (+24 more)

### Community 45 - "updatePage"
Cohesion: 0.15
Nodes (13): anyDaemonRunning(), installedBinary, latestRelease, updatePage, hashFile(), hashReader(), readInstalled(), probeSystem() (+5 more)

### Community 46 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.07
Nodes (15): recordTickCmd(), recordModel, diffModel, editorModel, listenForGuardEvents(), tickGuardStats(), renderBullets(), TestNoticeViewFitsAndPlacesOKLast() (+7 more)

### Community 48 - "runView"
Cohesion: 0.05
Nodes (18): daemonPage, daemonState, logGrid, logHook, logLine, runHost, runView, journalLine() (+10 more)

### Community 49 - "fadeScroll"
Cohesion: 0.05
Nodes (31): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+23 more)

### Community 50 - "time.Time"
Cohesion: 0.08
Nodes (24): LockedError, selfEventsFor(), sweepID, sweepMemo, EventType, parseEvents(), specFromForm(), specSummary() (+16 more)

### Community 51 - "protectedPage"
Cohesion: 0.05
Nodes (24): backendLink, configView, forwardForm, forwardView, installUI, protectedPage, recordReview, recordView (+16 more)

### Community 52 - "Client"
Cohesion: 0.08
Nodes (16): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+8 more)

### Community 53 - "highlight_test.go"
Cohesion: 0.13
Nodes (20): nextConfig(), fg(), Highlighter, lexerFor(), NewHighlighter(), shebangLexer(), editSession(), sampleGo() (+12 more)

### Community 54 - "fyne.io/fyne/v2/driver/desktop.Cursor"
Cohesion: 0.11
Nodes (12): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), newButtonCursor(), newCheckGroup(), newRadioGroup() (+4 more)

### Community 55 - "applier"
Cohesion: 0.15
Nodes (5): applier, ledgerOps, PlannedBinary, Store, vaultOps

### Community 56 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 57 - "Serve"
Cohesion: 0.08
Nodes (29): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop() (+21 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (64): engine, portOwner, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), confineConnect() (+56 more)

### Community 59 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 60 - "installPage"
Cohesion: 0.09
Nodes (15): actionGate, binSource, installAction, InstallMode, installPage, InstallRequest, Options, uninstallCard (+7 more)

### Community 61 - "deploy"
Cohesion: 0.15
Nodes (23): applyUpdate(), archAssets(), assetsFor(), checkUpdatePreconditions(), deploy(), deployGUI(), downloadAndVerify(), fetchGUI() (+15 more)

### Community 62 - "Ledger"
Cohesion: 0.12
Nodes (4): Ledger, Pending, Source, fakeLedger

### Community 63 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.11
Nodes (24): placedWalk, treeJudge, OnBtrfs(), SuperblockDev(), TestSuperblockDevMatchesStatOffBtrfs(), confineBelow(), inUserHome(), vouchedFd() (+16 more)

### Community 64 - "daemon.conf reference"
Cohesion: 0.10
Nodes (21): A complete example, `allow_lib`, Binary lines and events, daemon.conf reference, `[electron_apps]`, File format, `[inspectors]`, Installer placeholder (+13 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 67 - "wire.go"
Cohesion: 0.26
Nodes (9): decodeFrame(), distinctDictKeys(), header(), sameDict(), sameEntries(), sameList(), signatureField(), skipHeaderValue() (+1 more)

### Community 68 - "trustManager"
Cohesion: 0.09
Nodes (24): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree() (+16 more)

### Community 69 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 70 - "Vault"
Cohesion: 0.11
Nodes (25): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), deprovisionKind, SectionScan, isFileVaultCiphertext(), isRegularFileTarget() (+17 more)

### Community 71 - "CandidateDir"
Cohesion: 0.05
Nodes (42): symlinkTargetGlob(), collectDiffAdditions(), TestUncoveredCandidates(), uncoveredCandidates(), groupCandidates(), pickDirectories(), pickFromCandidates(), TestGroupCandidates() (+34 more)

### Community 72 - "changelog.go"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "FileTarget"
Cohesion: 0.06
Nodes (24): netGuardTab, policyEditor, policyHooks, policyTab, guardPill(), newPolicyTab(), verbFor(), netWarnings() (+16 more)

### Community 75 - "FileEvent"
Cohesion: 0.09
Nodes (9): snapMonitor, BpfEvent, eventUnitTest, FileEvent, Cstr(), DecodeBpfEvent(), MonitorRepository, fakeSource (+1 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "NetGuard"
Cohesion: 0.08
Nodes (14): lsmHook, PumpRingbuf(), reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs() (+6 more)

### Community 78 - "liveSession"
Cohesion: 0.14
Nodes (11): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), chooseResources(), eventsLabel(), forwardBinaries(), holdForward(), normalizedEvents() (+3 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 81 - "safeio.go"
Cohesion: 0.17
Nodes (18): dupFD(), installSSHAgent(), installSSHAgentEnv(), readAtNoFollow(), TestInstallSSHAgentRefusesSymlinkedUnitDir(), upsertUnitAt(), writeAtNoFollow(), TestWriteUserFileRefusesSymlinks() (+10 more)

### Community 82 - "watchSet"
Cohesion: 0.14
Nodes (13): addSystemBinary(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify, watchDir (+5 more)

### Community 83 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 84 - "fakeGuardRepo"
Cohesion: 0.05
Nodes (25): readTempJournal(), confirmUserPlaced(), forwardCheck(), binaryStatOf(), TempBinary, TemporaryGrant, TempRule, loadPinnedExeMaps() (+17 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.22
Nodes (14): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+6 more)

### Community 87 - "loginServer"
Cohesion: 0.09
Nodes (19): markReady(), markReadyWithLogin(), notifySystemdReady(), hasSentinel(), lockParent(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel() (+11 more)

### Community 88 - "compileGlobNames"
Cohesion: 0.28
Nodes (5): GlobChild, compileGlobNames(), globText(), resolveChildren(), TestCompileGlobNames_RepeatedPatternOrsBits()

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LaunchRule"
Cohesion: 0.09
Nodes (31): allowedEnv(), binaryLaunchRules(), launchRulesByPath(), TestAllowedEnvKeepsOnlyRootPlaced(), TestLaunchRulesByPath(), ShortOpts, argvEntry, applyLaunchRule() (+23 more)

### Community 92 - "User"
Cohesion: 0.05
Nodes (64): entryWriters(), reserveChain(), TestGlobBuilder_LibBitsPerEntry(), underResource(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor() (+56 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 95 - "dialLiveSession"
Cohesion: 0.07
Nodes (35): auditAfterEdit(), auditFindings(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), applyEdit(), checkAndApply() (+27 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "time.Duration"
Cohesion: 0.13
Nodes (15): procStats, statsSample, duration(), formatClock(), cpuTime(), humanBytes(), rssBytes(), New() (+7 more)

### Community 100 - "recordFlow"
Cohesion: 0.06
Nodes (39): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), lines(), offerBackupDeletion(), pickBinaries() (+31 more)

### Community 101 - "journal"
Cohesion: 0.15
Nodes (5): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, journal

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "ptrace_seize.c"
Cohesion: 0.20
Nodes (11): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file(), main(), peek_buffer(), run_tracer() (+3 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "vaultfs.go"
Cohesion: 0.09
Nodes (26): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode() (+18 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "planUpdaters"
Cohesion: 0.13
Nodes (20): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+12 more)

### Community 110 - "newTestTarget"
Cohesion: 0.22
Nodes (17): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+9 more)

### Community 111 - "Gate"
Cohesion: 0.16
Nodes (7): userGate, connectionGone(), Gate, Listen(), TestConnectionGoneOnlyFromTheBus(), Bus, Judge

### Community 112 - "UnifiedDiff"
Cohesion: 0.14
Nodes (18): refreshCatalogConfig(), installConfig(), confirmUnit(), ConfirmOverwrite(), ShowDiff(), TestUnifiedDiff(), TestUnifiedDiffIdentical(), UnifiedDiff() (+10 more)

### Community 113 - "NetworkMonitor"
Cohesion: 0.15
Nodes (7): BinariesSummary(), TestBinariesSummary(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey(), NetworkMonitor

### Community 115 - "reloadOnce"
Cohesion: 0.11
Nodes (22): newControlManager(), startControlManager(), backingDeviceUnion(), buildGuards(), makeReloadHandler(), reloadOnce(), reloadSlotsNeeded(), runDaemon() (+14 more)

### Community 116 - "Row"
Cohesion: 0.29
Nodes (10): candidates(), confirmAll(), describe(), record(), run(), Given(), Row, Inspect() (+2 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "catalogPatterns"
Cohesion: 0.50
Nodes (4): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), catalogPatterns(), inHome()

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "inet_connect.c"
Cohesion: 0.25
Nodes (17): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+9 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "validateConfigText"
Cohesion: 0.13
Nodes (16): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), TestAskEncryptionSkipsNeedEncryptionFalse() (+8 more)

### Community 127 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "SelectOrphans"
Cohesion: 0.16
Nodes (20): modifiedContextWithSource(), appProtectorSet(), CleanOrphanedMetadata(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors() (+12 more)

### Community 130 - ".exportTable"
Cohesion: 0.10
Nodes (13): binaryList, newBinaryList(), handCursors(), handCursorsOn(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog() (+5 more)

### Community 131 - "github.com/spf13/cobra.Command"
Cohesion: 0.09
Nodes (34): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+26 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "guardModel"
Cohesion: 0.06
Nodes (33): listenForDaemonEvents(), tickDaemonStats(), formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts(), formatDecision(), appendCapped(), clipToWidth() (+25 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "session"
Cohesion: 0.16
Nodes (5): peerCred(), readFrame(), queued, senderQueue, session

### Community 143 - "KernelDev"
Cohesion: 0.06
Nodes (43): BusClient, fdStat, rootSlot, BinaryStat, physicalParent(), slotOf(), statAt(), statFD() (+35 more)

### Community 144 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 145 - "Resource"
Cohesion: 0.05
Nodes (30): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), systemPatterns(), systemRule(), formatElapsed(), orNone() (+22 more)

### Community 146 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.15
Nodes (15): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle() (+7 more)

### Community 147 - "maintainerKeeper"
Cohesion: 0.18
Nodes (12): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp (+4 more)

### Community 148 - "btrfsLayoutFrom"
Cohesion: 0.17
Nodes (9): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), btrfsTypes(), TestBtrfsLayoutFromNestedMembers() (+1 more)

### Community 149 - "Recorder"
Cohesion: 0.16
Nodes (15): Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start() (+7 more)

### Community 150 - "SampleContent"
Cohesion: 0.17
Nodes (16): TestDaemonUnitWithMetadataOutput(), installFileAs(), installServices(), installUnit(), upsertFile(), isInstallerSSHAgentUnit(), TestIsInstallerSSHAgentUnit(), TestDaemonUnitCmdPrintsBundledUnit() (+8 more)

### Community 152 - "Backup"
Cohesion: 0.36
Nodes (11): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), Delete(), Find(), Backup, Restore() (+3 more)

### Community 154 - ".ExpandLibDirWriters"
Cohesion: 0.15
Nodes (13): homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelistSymlinkEscape(), TestSteamLibDirsCoverProtonCodeOnly(), TestSteamLibDirsExcludeScoutRuntime(), TestSteamLibDirWritersCoverItsOwnBinaries(), boundedGlob() (+5 more)

### Community 155 - "go_pkg_strings"
Cohesion: 0.05
Nodes (51): atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync(), writeWithin() (+43 more)

### Community 156 - "copytree.go"
Cohesion: 0.21
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 157 - "verdict"
Cohesion: 0.33
Nodes (8): arg(), attributesOf(), classifyCollection(), classifyProps(), classifyService(), withSchema(), verdict, verdictKind

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "NetGuardEvent"
Cohesion: 0.18
Nodes (6): snapNetGuard, NetGuardEvent, protoString(), netGuardEventLine, netGuardEventMsg, netGuardModel

### Community 161 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 162 - "backend.go"
Cohesion: 0.11
Nodes (28): asRoot(), checkInstalledMatchesSource(), ensureInstalledBinary(), fileSHA256(), prepareInstallation(), runInstall(), sameContent(), TestEnsureInstalledBinary() (+20 more)

### Community 164 - "serverAuth"
Cohesion: 0.32
Nodes (5): TestClientAuth(), TestSASL(), clientAuth(), serverAuth(), saslServer

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (46): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+38 more)

### Community 166 - "shell"
Cohesion: 0.04
Nodes (59): fakeLink, pageDef, pageID, reel, resizeWatch, shell, sizedDialog, snapEngines (+51 more)

### Community 167 - "classify"
Cohesion: 0.20
Nodes (12): call(), dictFrame(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyItemTypeIsSchema(), TestClassifyOpenCallsExact(), TestDuplicateDictKeyRefused(), TestReadFrame() (+4 more)

### Community 168 - "sources_test.go"
Cohesion: 0.32
Nodes (5): newSources(), TestArchiveMembersIgnoresNonArchives(), TestSourcesHoldFilesAndArchiveMembers(), TestSourcesRefuseOtherFilesystem(), writeTar()

### Community 169 - "os.FileInfo"
Cohesion: 0.19
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 170 - "searchDirs"
Cohesion: 0.18
Nodes (10): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths(), resolveSoname() (+2 more)

### Community 171 - "sources"
Cohesion: 0.24
Nodes (6): TestGeneratedLibraryBlockRoundTrips(), pathByte(), pathRecord(), TestPathRecord(), libDirDev(), sources

### Community 172 - "Limitations"
Cohesion: 0.17
Nodes (12): Already-running processes, Binaries you can write, Commands that stop working, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root (+4 more)

### Community 173 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 174 - "plan_test.go"
Cohesion: 0.27
Nodes (14): isWhiteout(), opaque(), scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles() (+6 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "editauth_test.go"
Cohesion: 0.10
Nodes (27): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+19 more)

### Community 178 - "mustSubkey"
Cohesion: 0.21
Nodes (13): deriveFileVaultSubkey(), newFileVaultAEAD(), openFileVault(), sealFileVault(), sealFileVaultWithMasterKey(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape() (+5 more)

### Community 179 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 181 - "New"
Cohesion: 0.13
Nodes (33): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+25 more)

### Community 182 - "TrustGuard"
Cohesion: 0.06
Nodes (21): statFunc, tempGrant, trustHook, trustRows, ensureBtrfsLayout(), setSupersedeTrust(), liveTrust(), resolveBits() (+13 more)

### Community 186 - "Binary trust"
Cohesion: 0.20
Nodes (10): Binary ledger, Binary replacement, Binary trust, Browser profiles, Hash cache, Launch rules, Launch scan, Multicall binaries (+2 more)

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 191 - "os.File"
Cohesion: 0.04
Nodes (72): scriptIdentity(), checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir() (+64 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "run.go"
Cohesion: 0.36
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 194 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "ensurePortAnchors"
Cohesion: 0.29
Nodes (8): ensureAnchor(), ensurePortAnchors(), ensureRootDir(), ephemeralRange(), parseReservedPorts(), readReservedPorts(), TestParseReservedPorts(), ensureSecretAnchors()

### Community 202 - "auditGroupCoverage"
Cohesion: 0.22
Nodes (11): auditEntry(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+3 more)

### Community 204 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 205 - "SystemTrusted"
Cohesion: 0.33
Nodes (7): rootOwnedInode(), rootOwnedSafe(), SystemTrusted(), systemTool(), archiveMembers(), decompressed(), tarMembers()

### Community 206 - "monitor"
Cohesion: 0.07
Nodes (33): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+25 more)

### Community 209 - "absPath"
Cohesion: 0.05
Nodes (13): IntegrationSuite, IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, launchCase (+5 more)

### Community 210 - ".admitGated"
Cohesion: 0.24
Nodes (3): placeableAttributes(), TestLearnKeyringFromFirstReply(), TestPlaceableAttributes()

### Community 211 - "runMonitor"
Cohesion: 0.15
Nodes (16): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag(), pathWithinDir(), ResolveTargets(), validatePair() (+8 more)

### Community 212 - "catalogRefresher"
Cohesion: 0.14
Nodes (9): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestDebouncerPendingCoversWaitAndRun(), catalogRefresher (+1 more)

### Community 213 - "update_test.go"
Cohesion: 0.09
Nodes (22): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+14 more)

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "runNetworkGuard"
Cohesion: 0.08
Nodes (33): DrainUntilSignal(), CheckBPFLSM(), CheckEBPF(), ParseNetEventsFlag(), parseTypesFlag(), newPinGeneration(), prepareDaemonStart(), runBPFCheck() (+25 more)

### Community 218 - "runGenKey"
Cohesion: 0.29
Nodes (8): confirmMasterKeyOverwrite(), runGenKey(), runOneShotMode(), ensureMasterKey(), checkKeyLen(), GenerateMasterKey(), MasterKeyExists(), readOrCreateKey()

### Community 219 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 221 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 223 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 224 - "batchInodeRows"
Cohesion: 0.27
Nodes (9): inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches(), countInodeRows() (+1 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "decryptStep"
Cohesion: 0.22
Nodes (8): askDecrypt(), cleanOrphanedMetadata(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty(), VerifyEncryptedKeys(), vaultKeyCheck

### Community 227 - "bintrust_test.go"
Cohesion: 0.40
Nodes (9): ConfirmGiven(), ConfirmPending(), stillHolds(), openLedger(), TestConfirmGivenNeedsTheReviewedContent(), TestConfirmPendingNeverOpensALineTheLedgerDoesNotRefuse(), TestConfirmPendingOnlyWhatWasRefusedAndIsStillThere(), TestPendingRowsSkipChangedContent() (+1 more)

### Community 228 - "launch_probe.c"
Cohesion: 0.46
Nodes (6): anchor(), copy(), is_verb(), main(), read_once(), usage()

### Community 231 - "procstats.go"
Cohesion: 0.38
Nodes (8): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic()

### Community 233 - "secretJudge"
Cohesion: 0.15
Nodes (7): placedKey(), secretJudge, SecretService, LiveGuardAt(), Caller, Role, Scope

### Community 234 - "mutedText"
Cohesion: 0.08
Nodes (40): checkGrid, filesPage, handCheck, homePage, multiAnswer, recordForm, newBinaryRow(), newCheck() (+32 more)

### Community 235 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.13
Nodes (31): class, pos, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+23 more)

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), userMountable() (+2 more)

### Community 238 - "GenerateConf"
Cohesion: 0.05
Nodes (77): appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), editConfig(), libraryBlocksFromCandidates(), runConfigEditor(), sectionsFromCandidates(), sharedSection(), discordCatalogEntry() (+69 more)

### Community 241 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 242 - "linkView"
Cohesion: 0.36
Nodes (4): linkTarget(), entryKind, libPlan, linkView

### Community 243 - "bpfLsmListed"
Cohesion: 0.40
Nodes (6): bpfLsmListed(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 244 - "downloadFile"
Cohesion: 0.32
Nodes (6): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag()

### Community 245 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 246 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 247 - "ShortRefused"
Cohesion: 0.33
Nodes (5): ShortRefused(), TestSSHConfigFileRefusedInEveryCluster(), shortMask(), shortHit(), TestShortMaskRefusesNonLetters()

### Community 248 - "run"
Cohesion: 0.36
Nodes (8): InstalledGUIVersion(), InstalledVersion(), readInstalledVersion(), run(), runUpdate(), syncGUI(), TestReadInstalledVersion(), GUIInstalled()

### Community 249 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 252 - "CleanupStalePins"
Cohesion: 0.33
Nodes (6): CleanupStalePins(), foreignPinnedLink(), genOfPinFile(), TestCleanupStalePins(), TestCleanupStalePinsMissingBase(), TestGenOfPinFile()

### Community 255 - "monitor.md"
Cohesion: 0.50
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 256 - "network monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 257 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 259 - "noteReloadedSettings"
Cohesion: 0.40
Nodes (5): noteReloadedSettings(), ConfigureInodeCapacity(), InodeOccupancy(), TestConfigureInodeCapacityBounds(), WarnInodeOccupancy()

### Community 260 - "SecretMatch"
Cohesion: 0.40
Nodes (4): SecretScope, SecretMatch, parseSecretMatch(), TestSecretMatch()

### Community 261 - "`[launch_rules]`"
Cohesion: 0.40
Nodes (5): `allow_env`, `deny_args`, `warn_args`, `deny_env`, `warn_env`, `[launch_rules]`, refuses()

### Community 264 - "parsePasswd"
Cohesion: 0.50
Nodes (4): parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed()

### Community 267 - "daemonconfig.go"
Cohesion: 0.07
Nodes (74): runtimePathsFor(), watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyBinaryLine(), applyBlockLine(), applyConfigLine() (+66 more)

## Knowledge Gaps
- **34 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck` (+29 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 691 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **52 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `guardUnitTest`, `GuardEvent`, `ConfinedEntry`, `buildOneGuard`, `DaemonEvent`, `KernelDev`, `Resource`, `BinaryEntry`, `sync.Mutex`, `time.Time`, `TrustGuard`, `GuardInodeKey`, `os.File`, `Vault`, `batchInodeRows`, `time.Duration`, `secretJudge`, `reloadOnce`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _34 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14829594477998273 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `parseDaemonEvents`, `IntegrationSuite`, `IntegrationSuite`, `.netIntegrityContainer`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `go_pkg_strings`, `IntegrationSuite`?**
  _High betweenness centrality (0.033) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.015280223764632757 - nodes in this community are weakly interconnected._
- **Why does `monitor()` connect `monitor` to `time.Duration`, `secret_probe/main.go`, `Installation`, `app-listener documentation`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09441441441441441 - nodes in this community are weakly interconnected._