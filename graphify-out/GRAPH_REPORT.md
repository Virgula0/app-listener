# Graph Report - app-listener  (2026-10-18)

## Corpus Check
- 471 files · ~1,692,535 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6177 nodes · 26264 edges · 219 communities (178 shown, 41 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2479 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `7cafbddf`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- rcUser
- guardUnitTest
- mkdirs
- changelog.go
- policyEditor
- protectedPage
- github.com/charmbracelet/bubbletea.Cmd
- go_pkg_github_com_testcontainers_testcontainers_go
- editView
- IntegrationSuite
- guard_trust.bpf.c
- IntegrationSuite
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- sync.Mutex
- guard.bpf.c
- NetEventType
- time.Time
- .FilterExistingWhitelist
- go_pkg_path_filepath
- ui.go
- serveWebSocket
- IntegrationSuite
- netModel
- picker
- github.com/spf13/cobra.Command
- IntegrationSuite
- backend
- controlServer
- fileedit_model_test.go
- os.File
- IntegrationSuite
- IntegrationSuite
- monitorUnitTest
- .objs
- GuardEvent
- .posAt
- stdio
- installcert.go
- fadeScroll
- .showSized
- GenerateConf
- Client
- Config
- navItem
- editauth_test.go
- run
- icon
- GuardInodeKey
- New
- User
- daemon_multicall_test.go
- exe_supersede.h
- Network
- binaryList
- bpfstats/main.go
- wm.go
- runMonitor
- fakeGuardRepo
- KernelDev
- fcntl
- Parse
- Guard
- FileTarget
- LibraryClosure
- catalogRefresher
- RefreshSection
- audit.go
- process_vm_readv.c
- Monitor
- journal
- planUpdaters
- updatePage
- compileGlobNames
- highlight_test.go
- github.com/charmbracelet/bubbles/textarea.Model
- btrfsLayoutFrom
- liveSession
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- DaemonEvent
- changelog_test.go
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- newTestTarget
- fileedit_symlink_test.go
- update/update.go
- Well-known bypass classes
- rootHandle
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- tempGrantSetup
- Build & sign release assets (reusable workflow)
- installPage
- serve.go
- mountTable
- NetGuardEvent
- buildTrustedSet
- IntegrationSuite
- IntegrationSuite
- net_tester/main.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- integrationtests/networkguard_test.go
- IntegrationSuite
- WebSocket /ws client
- launch_probe.c
- install.sh
- mc_tag.h
- session_test.go
- make test-integration (rootful Docker suite)
- hashFileShared
- supersede_probe.c
- go_pkg_strings
- trace-app-libs.sh
- IntegrationSuite
- RawTarget
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- time.Duration
- io.Writer
- hashCache
- runDaemon
- github.com/stretchr/testify/suite.Suite
- appTheme
- uninstall_system.go
- makeReloadHandler
- IntegrationSuite
- serveHandler
- ConfinedEntry
- EventType
- layout.go
- copyRegularAt
- NewMonitor
- Installation
- backend_test.go
- graphify-refresh.sh
- Vault
- inodeChain
- configedit_test.go
- pinstate.go
- feed.go
- shell
- runView
- .dispatch
- AddServeFlags
- Finder
- ParseSectionWhitelist
- Limitations
- Troubleshooting
- IntegrationSuite
- fakeVault
- Daemon configuration
- NewDumpHook
- NewEventFanout
- install/install.go
- Compatibility
- How it works
- parseResize
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- validateFlags
- root.go
- edit-protected
- monitor: see what touches a path
- kernel-prof.sh
- fakeDaemon
- pathCache
- IntegrationSuite
- formatEntry
- inode_dev.h
- runNetworkGuard
- parseGuardEvent
- Desktop GUI
- IntegrationSuite
- Ledger
- .runEditorHarness
- .eventTypeAllowed
- ParseMode
- TestCleanRel
- bleedLayout
- Vault
- SanitizeText
- TestMain
- partitionEncryptionRoots

## God Nodes (most connected - your core abstractions)
1. `Guard` - 118 edges
2. `GuardInodeKey` - 99 edges
3. `Config` - 95 edges
4. `Load()` - 89 edges
5. `shell` - 87 edges
6. `IntegrationSuite` - 82 edges
7. `absPath()` - 81 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 69 edges
10. `IntegrationSuite` - 68 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `resolveAssets()` --calls--> `asset()`  [INFERRED]
  cmd/functions/update/update.go → internal/gui/icons.go

## Import Cycles
- None detected.

## Communities (219 total, 41 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.16
Nodes (26): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey, launchModel() (+18 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (128): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), awaitStartupOrSignal(), loadDaemonConfig(), setConfinedHomes(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod() (+120 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (10): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote() (+2 more)

### Community 3 - "Vault"
Cohesion: 0.07
Nodes (31): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), deprovisionKind (+23 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.08
Nodes (73): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+65 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (32): emit_event(), emit_event_kern(), fill_identity(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice() (+24 more)

### Community 7 - "rcUser"
Cohesion: 0.13
Nodes (34): setupBunTmpdir(), rcTarget, BunTmpDir(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), EnsureBunTmpDir(), ensureRCBlock() (+26 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (20): addErrKind, ReplacementCheck, BinariesSummary(), classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+12 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 11 - "policyEditor"
Cohesion: 0.14
Nodes (11): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+3 more)

### Community 12 - "protectedPage"
Cohesion: 0.07
Nodes (13): backendLink, configView, forwardForm, forwardView, protectedPage, Event, Hello, showInformation() (+5 more)

### Community 13 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (24): diffModel, editorModel, listenForDaemonEvents(), NewDaemonModel(), tickDaemonStats(), syncViewport(), protoString(), renderBullets() (+16 more)

### Community 14 - "go_pkg_github_com_testcontainers_testcontainers_go"
Cohesion: 0.09
Nodes (8): attrOpCase, IntegrationSuite, inNosuidNamespace(), guardEventTypes(), guardExploitTest, launchCase, pooledContainer, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 15 - "editView"
Cohesion: 0.11
Nodes (9): codeEntry, editView, handButton, Entry, selectable(), fileSize(), joinRel(), newCodeEntry() (+1 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.09
Nodes (8): IntegrationSuite, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, absPath()

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (58): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+50 more)

### Community 18 - "IntegrationSuite"
Cohesion: 0.24
Nodes (6): netMonitorEvent, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (11): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, fileEditMode (+3 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.20
Nodes (44): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+36 more)

### Community 22 - "Resource"
Cohesion: 0.08
Nodes (9): Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots(), unlockRoots(), vaultOpForGuard() (+1 more)

### Community 23 - "auth.go"
Cohesion: 0.12
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "sync.Mutex"
Cohesion: 0.06
Nodes (20): refusedReplacements, trustHook, trustRows, ensureBtrfsLayout(), deleteInoKeys(), dropSupersededLocked(), PruneSuperseded(), setSupersedeTrust() (+12 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (72): inode_dev(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+64 more)

### Community 26 - "NetEventType"
Cohesion: 0.06
Nodes (20): snapNetMonitor, NetBpfEvent, NetEvent, NetEventType, eventsOrAll(), newNetEventChecks(), FormatAddr(), NetEventTypes() (+12 more)

### Community 27 - "time.Time"
Cohesion: 0.09
Nodes (21): LockedError, addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), fakeInotify (+13 more)

### Community 28 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 29 - "go_pkg_path_filepath"
Cohesion: 0.33
Nodes (7): isAbsoluteCandidate(), TestCatalogSanity(), parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed(), hint

### Community 30 - "ui.go"
Cohesion: 0.07
Nodes (33): nextConfig(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), TestConfirmUpdateRoutesToActiveUI() (+25 more)

### Community 31 - "serveWebSocket"
Cohesion: 0.21
Nodes (9): configureWebSocket(), readResizeLoop(), sameOrigin(), serveWebSocket(), TestSameOrigin(), TestSnapshotHubSingleViewerAndCoalescing(), writeSnapshot(), writeSnapshotLoop() (+1 more)

### Community 33 - "netModel"
Cohesion: 0.06
Nodes (30): runNetworkMonitor(), runTUI(), Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace() (+22 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (30): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+22 more)

### Community 35 - "github.com/spf13/cobra.Command"
Cohesion: 0.27
Nodes (11): credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestSetupDumpLogRefusesSymlink(), effectiveVerbose() (+3 more)

### Community 37 - "backend"
Cohesion: 0.08
Nodes (19): checkSecret(), TestLiveSessionWatchAndPinger(), displayEntries(), holdForward(), pingInterval(), revertBunLaunchers(), File, backend (+11 more)

### Community 38 - "controlServer"
Cohesion: 0.10
Nodes (22): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest(), readPathLines() (+14 more)

### Community 39 - "fileedit_model_test.go"
Cohesion: 0.18
Nodes (27): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+19 more)

### Community 40 - "os.File"
Cohesion: 0.04
Nodes (60): VettedInode, Multicall, MulticallKind, placedWalk, ExeKey(), IsMulticallRefusal(), multicallFamily(), MulticallRefusal() (+52 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.11
Nodes (5): IntegrationSuite, exploitTest, IntegrationSuite, newEventTypes(), IntegrationSuite

### Community 44 - ".objs"
Cohesion: 0.08
Nodes (10): exeRow, tempGrant, tempMaskOp, tempRow, eventMask(), supersededMark(), liveTrust(), loadPinnedExeMaps() (+2 more)

### Community 45 - "GuardEvent"
Cohesion: 0.09
Nodes (20): holdHeadless(), HeadlessLine(), resolveGuardConfig(), runGuard(), runGuardHeadless(), runGuardTUI(), snapGuard, GuardEvent (+12 more)

### Community 46 - ".posAt"
Cohesion: 0.13
Nodes (12): lineStyles, viewportState, fg(), baseStyle(), gutterLabel(), lineLen(), markStyle(), setGutter() (+4 more)

### Community 47 - "stdio"
Cohesion: 0.07
Nodes (4): put(), wait_then_read(), hold_forever(), main()

### Community 48 - "installcert.go"
Cohesion: 0.06
Nodes (44): dupFD(), installFile(), installFileAs(), installServices(), installSSHAgent(), installSSHAgentEnv(), installUnit(), readAtNoFollow() (+36 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (26): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+18 more)

### Community 50 - ".showSized"
Cohesion: 0.12
Nodes (10): installUI, resizeWatch, sizedDialog, scrollable(), fadeScrollOf(), newFadeScroll(), ask(), diffText() (+2 more)

### Community 51 - "GenerateConf"
Cohesion: 0.06
Nodes (51): appendSectionsAndEdit(), TestDiffMergePreservesExistingAndParses(), libraryBlocksFromCandidates(), sectionsFromCandidates(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH() (+43 more)

### Community 52 - "Client"
Cohesion: 0.08
Nodes (18): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+10 more)

### Community 53 - "Config"
Cohesion: 0.06
Nodes (63): newBinaryVetter(), openBinaryVetter(), sameConfig(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), startTrustGuard(), trustStartupError(), collectDiffAdditions() (+55 more)

### Community 54 - "navItem"
Cohesion: 0.03
Nodes (22): buttonCursor, clickPanel, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabBar (+14 more)

### Community 55 - "editauth_test.go"
Cohesion: 0.12
Nodes (26): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+18 more)

### Community 56 - "run"
Cohesion: 0.08
Nodes (32): deletePostBackups(), listBackups(), restoreBackups(), askDecrypt(), cleanOrphanedMetadata(), confirmKeyDeletion(), decryptDirectories(), decryptStep() (+24 more)

### Community 57 - "icon"
Cohesion: 0.09
Nodes (30): checkGrid, filesPage, homePage, multiAnswer, newBinaryList(), newBinaryRow(), TestAdaptiveGridDropsColumns(), newFadeVScroll() (+22 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.04
Nodes (38): engine, resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded(), appendInodeRow(), batchInodeRows() (+30 more)

### Community 59 - "New"
Cohesion: 0.08
Nodes (47): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+39 more)

### Community 60 - "User"
Cohesion: 0.05
Nodes (45): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), bunEntryConfigured(), bunLaunchersFor(), selectAndEditConfig() (+37 more)

### Community 61 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.21
Nodes (15): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+7 more)

### Community 63 - "Network"
Cohesion: 0.13
Nodes (9): netGuardPill(), guardRow(), clonePolicy(), NetExtra, NetPolicy, Network, verifyHashes(), exeCached (+1 more)

### Community 64 - "binaryList"
Cohesion: 0.13
Nodes (10): binaryList, handCursors(), handCursorsOn(), showConfirm(), showCustomConfirm(), showDialog(), showErrorDialog(), showForm() (+2 more)

### Community 65 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (26): invokingUID(), detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket() (+18 more)

### Community 67 - "runMonitor"
Cohesion: 0.20
Nodes (8): MakeDisplayPaths(), ParseEventsFlag(), validateForwardRule(), runHeadless(), runMonitor(), runTUI(), NewModel(), MonitorUseCase

### Community 69 - "KernelDev"
Cohesion: 0.13
Nodes (22): BinaryStat, isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs(), binaryStatOf() (+14 more)

### Community 70 - "fcntl"
Cohesion: 0.09
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 71 - "Parse"
Cohesion: 0.07
Nodes (33): auditAfterEdit(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), applyEdit(), checkAndApply(), editAgain() (+25 more)

### Community 72 - "Guard"
Cohesion: 0.06
Nodes (28): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), deferredBinary, keptEntries, BinaryEntry, TestWithChmodDropWriteNeedsReadOnly(), checkBtrfsKeys() (+20 more)

### Community 73 - "FileTarget"
Cohesion: 0.15
Nodes (9): guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps() (+1 more)

### Community 74 - "LibraryClosure"
Cohesion: 0.07
Nodes (29): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+21 more)

### Community 76 - "catalogRefresher"
Cohesion: 0.14
Nodes (9): drain(), newCatalogRefresher(), readEvents(), swapFile(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), startCatalogRefresh(), catalogRefresher (+1 more)

### Community 77 - "RefreshSection"
Cohesion: 0.10
Nodes (28): logRefreshChange(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange(), refreshCatalogConfig() (+20 more)

### Community 78 - "audit.go"
Cohesion: 0.12
Nodes (26): auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+18 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.18
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 81 - "journal"
Cohesion: 0.14
Nodes (7): fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal, NetGuardConfig

### Community 82 - "planUpdaters"
Cohesion: 0.15
Nodes (17): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+9 more)

### Community 83 - "updatePage"
Cohesion: 0.23
Nodes (4): latestRelease, updatePage, boolIndex(), isPreRelease()

### Community 84 - "compileGlobNames"
Cohesion: 0.17
Nodes (9): GlobChild, statFunc, compileGlobNames(), globText(), resolveBits(), resolveChildren(), statKey(), TestCompileGlobNames_RepeatedPatternOrsBits() (+1 more)

### Community 85 - "highlight_test.go"
Cohesion: 0.17
Nodes (14): lexerFor(), NewHighlighter(), shebangLexer(), sampleGo(), TestDaemonConfTokens(), TestGutterKeepsTextColumnAligned(), TestHighlighterColorsTokens(), TestHighlighterFollowsCursorAlone() (+6 more)

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.10
Nodes (43): class, pos, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+35 more)

### Community 87 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 88 - "liveSession"
Cohesion: 0.22
Nodes (10): parseErr(), TestLiveSessionConfigAndPut(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), forwardCheck(), normalizedEvents() (+2 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "DaemonEvent"
Cohesion: 0.11
Nodes (30): ServeConfig, UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI() (+22 more)

### Community 92 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.07
Nodes (21): NewMountVouch(), reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), codeVouched(), isTraced(), judgePasses() (+13 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "newTestTarget"
Cohesion: 0.21
Nodes (18): newExeResolver(), FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps() (+10 more)

### Community 100 - "fileedit_symlink_test.go"
Cohesion: 0.42
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 101 - "update/update.go"
Cohesion: 0.04
Nodes (87): installHooks(), staged(), changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), showChangelog(), applyUpdate() (+79 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "rootHandle"
Cohesion: 0.12
Nodes (9): rootHandle, openRoot(), TestNoteRootMountFollowsMountID(), TestReadOnlyWalkToleratesDanglingSymlinks(), FirstOnBtrfs(), OnBtrfs(), SuperblockDev(), SuperblockDevPath() (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "installPage"
Cohesion: 0.06
Nodes (23): actionGate, binaryRow, binSource, iconTile, installAction, installedBinary, InstallMode, installPage (+15 more)

### Community 110 - "serve.go"
Cohesion: 0.27
Nodes (10): isInteractiveTerminal(), newServeListener(), newServeServer(), relayServeSignals(), runServe(), Serve(), serveHTTP(), waitProgram() (+2 more)

### Community 111 - "mountTable"
Cohesion: 0.18
Nodes (9): MountInfo, mountinfoUnitTest, MountByID(), parseMountDev(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo() (+1 more)

### Community 112 - "NetGuardEvent"
Cohesion: 0.09
Nodes (19): runHeadless(), runHeadless(), snapNetGuard, ProtocolString(), NetGuardEvent, NetworkGuardRepository, NewGuardUseCase(), NewMonitorUseCase() (+11 more)

### Community 113 - "buildTrustedSet"
Cohesion: 0.13
Nodes (21): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), inspectorPaths(), TestInspectorPathsOnlyRootPlaced() (+13 more)

### Community 116 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "integrationtests/networkguard_test.go"
Cohesion: 0.27
Nodes (5): IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "launch_probe.c"
Cohesion: 0.53
Nodes (4): copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "session_test.go"
Cohesion: 0.19
Nodes (13): ancestorEnv(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue(), TestOwnedByChecksTypeAndOwner(), TestReadOwnedFileRefusesSymlink() (+5 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "hashFileShared"
Cohesion: 0.19
Nodes (10): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+2 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "go_pkg_strings"
Cohesion: 0.07
Nodes (19): releaseBuild(), Run(), User, stageItem, stageState, textAnswer, infraContainerPaths(), TestIntegrationSuite() (+11 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "RawTarget"
Cohesion: 0.27
Nodes (10): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+2 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.14
Nodes (15): TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+7 more)

### Community 143 - "time.Duration"
Cohesion: 0.15
Nodes (13): procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv(), TestSubDeliversMissedTickLate() (+5 more)

### Community 144 - "io.Writer"
Cohesion: 0.27
Nodes (9): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), fyneLogFilter, ReloadDaemonIfActive() (+1 more)

### Community 145 - "hashCache"
Cohesion: 0.19
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 146 - "runDaemon"
Cohesion: 0.08
Nodes (33): CheckBPFLSM(), CheckEBPF(), backingDeviceUnion(), buildConcurrency(), buildGuards(), catchLifecycleSignals(), confirmMasterKeyOverwrite(), newPinGeneration() (+25 more)

### Community 149 - "uninstall_system.go"
Cohesion: 0.06
Nodes (62): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeKeyAndEmptyDir(), removeMasterKey(), removeSSHAgentEnv(), removeSSHAgentUnit() (+54 more)

### Community 150 - "makeReloadHandler"
Cohesion: 0.15
Nodes (11): newControlManager(), startControlManager(), makeReloadHandler(), markReady(), noteReloadedSettings(), notifySystemdReady(), controlManager, ConfigureInodeCapacity() (+3 more)

### Community 152 - "serveHandler"
Cohesion: 0.31
Nodes (8): basicAuth(), securityHeaders(), serveHandler(), TestBasicAuth(), TestSecurityHeadersAndExactRoot(), TestServeHandlerProtectsPageAndWebSocket(), TestServeWebSocketSnapshotsResizeAndSingleViewer(), validateHost()

### Community 153 - "ConfinedEntry"
Cohesion: 0.29
Nodes (18): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+10 more)

### Community 154 - "EventType"
Cohesion: 0.07
Nodes (17): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction() (+9 more)

### Community 155 - "layout.go"
Cohesion: 0.43
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 156 - "copyRegularAt"
Cohesion: 0.15
Nodes (18): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), CopyTree(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo() (+10 more)

### Community 157 - "NewMonitor"
Cohesion: 0.22
Nodes (3): NewMonitor(), newMountTable(), RealFileEngines()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "backend_test.go"
Cohesion: 0.22
Nodes (19): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+11 more)

### Community 161 - "Vault"
Cohesion: 0.09
Nodes (26): checkFreshDir(), openDirNoFollow(), renewDir(), entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent() (+18 more)

### Community 162 - "inodeChain"
Cohesion: 0.21
Nodes (16): fdStat, inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), physicalParent(), statFD(), keyOf() (+8 more)

### Community 163 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 164 - "pinstate.go"
Cohesion: 0.11
Nodes (33): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder() (+25 more)

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (46): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+38 more)

### Community 166 - "shell"
Cohesion: 0.06
Nodes (46): pageDef, pageID, reel, shell, snapEngines, snapNet, sysStatus, newLogo() (+38 more)

### Community 167 - "runView"
Cohesion: 0.07
Nodes (9): logGrid, logHook, logLine, runHost, runView, logRow(), logStyle(), newLogGrid() (+1 more)

### Community 168 - ".dispatch"
Cohesion: 0.22
Nodes (4): TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestDispatchWithoutOwnAllowedEvents(), commMatchesGuardedBinary()

### Community 169 - "AddServeFlags"
Cohesion: 0.29
Nodes (7): AddServeFlags(), TestParseServeFlagsRequiresCredentials(), init(), init(), init(), init(), init()

### Community 170 - "Finder"
Cohesion: 0.13
Nodes (20): findButton, span, findAll(), Finder, lowerRunes(), NewFinder(), newFindInput(), overlay() (+12 more)

### Community 171 - "ParseSectionWhitelist"
Cohesion: 0.29
Nodes (7): isLibDirectiveText(), ParseSectionHeaderPath(), unquotePath(), IsLibraryDirective(), parseLibDirLine(), ParseSectionWhitelist(), TestParseSectionWhitelist()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 175 - "fakeVault"
Cohesion: 0.42
Nodes (3): TestDaemonUseCaseReloadRefusedDuringShutdown(), TestDaemonUseCaseStopBlocksUntilAllLocked(), fakeVault

### Community 176 - "Daemon configuration"
Cohesion: 0.25
Nodes (8): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Several guarded trees, one vault, Watch sections

### Community 177 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 178 - "NewEventFanout"
Cohesion: 0.22
Nodes (5): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout

### Community 179 - "install/install.go"
Cohesion: 0.07
Nodes (51): askBunTmpdirUsers(), gatherPerUserSetup(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall() (+43 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "parseResize"
Cohesion: 0.40
Nodes (3): decodeEndsAtEOF(), parseResize(), TestParseResize()

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 194 - "validateFlags"
Cohesion: 0.33
Nodes (5): TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), TestUnchangedPerShared()

### Community 195 - "root.go"
Cohesion: 0.10
Nodes (16): Requested(), TestRequested(), atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget() (+8 more)

### Community 196 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 202 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 203 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 204 - "runNetworkGuard"
Cohesion: 0.15
Nodes (15): ParseNetEventsFlag(), computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), WithSuspectsFrom() (+7 more)

### Community 205 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.15
Nodes (13): Actions, Desktop GUI, File and folder pickers, Files, Home, Install, Network, Protected (+5 more)

### Community 209 - "Ledger"
Cohesion: 0.14
Nodes (4): applySchema(), Ledger, Pending, Source

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "ParseMode"
Cohesion: 0.50
Nodes (3): IsOctalMode(), ParseMode(), TestIsOctalMode()

### Community 216 - "SanitizeText"
Cohesion: 0.07
Nodes (37): multicallAdmissible(), ownerPaths(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), buildOneGuard(), eventFilterOptions(), systemPatterns() (+29 more)

### Community 217 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

## Knowledge Gaps
- **31 isolated node(s):** `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC`, `statonly stat/statx metadata leak`, `app-listener-gui best-effort build` (+26 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 567 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **41 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `hashFileShared`, `fakeDaemon`, `guardUnitTest`, `os.File`, `time.Time`, `rootHandle`, `.objs`, `.dispatch`, `GuardEvent`, `time.Duration`, `runDaemon`, `makeReloadHandler`, `SanitizeText`, `sync.Mutex`, `GuardInodeKey`, `DaemonEvent`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **What connects `IntegrationSuite`, `github.com/Virgula0/app-listener`, `open_by_handle_at POC` to the rest of the system?**
  _31 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.020934536201711776 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `IntegrationSuite` to `IntegrationSuite`, `go_pkg_strings`, `TestMain`, `IntegrationSuite`, `IntegrationSuite`, `go_pkg_github_com_testcontainers_testcontainers_go`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `integrationtests/networkguard_test.go`, `IntegrationSuite`, `daemon_multicall_test.go`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.08958130477117819 - nodes in this community are weakly interconnected._
- **Why does `IntegrationSuite` connect `IntegrationSuite` to `.newGuardTestContainer`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `Vault` be split into smaller, more focused modules?**
  _Cohesion score 0.06829573934837092 - nodes in this community are weakly interconnected._