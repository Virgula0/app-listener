# Graph Report - app-listener  (2026-11-07)

## Corpus Check
- 575 files · ~1,819,972 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7540 nodes · 34931 edges · 258 communities (213 shown, 45 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 3096 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `3c5eab35`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- go_pkg_strings
- testing.T
- IntegrationSuite
- SecretStoreNotes
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- .dispatch
- ConfinedEntry
- RefreshSection
- Guard
- buildOneGuard
- time.Time
- IntegrationSuite
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- op
- guard.bpf.c
- github.com/charmbracelet/bubbles/textarea.Model
- inode_dev.h
- wm.go
- rcUser
- editView
- Network
- CandidateDir
- controlServer
- picker
- secret_probe/main.go
- IntegrationSuite
- Message
- systemd.go
- NetEventType
- desktop.go
- UI
- fileedit_model_test.go
- monitorUnitTest
- pinstate.go
- cursor.go
- github.com/charmbracelet/bubbletea.Cmd
- stdio
- runView
- fadeScroll
- EventType
- shell
- editbackend_test.go
- changelog.go
- Gate
- .apply
- backend_test.go
- serve.go
- GuardInodeKey
- IntegrationSuite
- FileTarget
- runMonitor
- runDaemon
- golang.org/x/sys/unix.Stat_t
- gui/session.go
- unix_connect.c
- IntegrationSuite
- wire.go
- buildTrustedSet
- parseDaemonEvents
- Vault
- expandPlaceholders
- runNetworkMonitor
- fcntl
- policyEditor
- gui.go
- codecheck/main.go
- NetGuard
- liveSession
- process_vm_readv.c
- .filterEvent
- bufio.Reader
- daemon.conf reference
- daemon_multicall_test.go
- fakeGuardRepo
- bpfstats/main.go
- exe_supersede.h
- loginServer
- daemonPage
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LaunchRule
- User
- netinject.c
- README.md
- dialLiveSession
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- recordFlow
- newTestTarget
- Well-known bypass classes
- main_test.go
- jit_provenance.c
- ptrace_seize.c
- recordops.c
- vaultfs.go
- Build & sign release assets (reusable workflow)
- planUpdaters
- releasesSince
- run.go
- SuperblockDev
- validateConfigText
- IntegrationSuite
- catalogRefresher
- Row
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- NetGuardEvent
- IntegrationSuite
- WebSocket /ws client
- watchPattern
- inet_connect.c
- install.sh
- watchSet
- protected_test.go
- make test-integration (rootful Docker suite)
- SelectOrphans
- Recorder
- github.com/spf13/cobra.Command
- trace-app-libs.sh
- IntegrationSuite
- ValidateTargets
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- session
- KernelDev
- Monitor
- Ledger
- TestMonitorUseCaseLifecycle
- maintainerKeeper
- btrfsLayoutFrom
- github.com/stretchr/testify/suite.Suite
- decryptStep
- fakeStore
- uninstallAll
- IntegrationSuite
- maintainer/plan.go
- guardLSMHooks
- io/fs.FileInfo
- classify
- Installation
- netGuardModel
- graphify-refresh.sh
- layout.go
- update_test.go
- SystemTrusted
- mc_tag.h
- icon
- TestSnapshots
- secureResources
- VouchedDevs
- sources_test.go
- `[watch <path>]`
- os.FileInfo
- Limitations
- .keyringObject
- plan_test.go
- net_tester/main.go
- Open
- sources
- tempGrantSetup
- IntegrationSuite
- filevault_test.go
- TrustGuard
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- IntegrationSuite
- os.File
- Self-updating apps and the catalog
- sync.Mutex
- WithBottomBar
- edit-protected
- inspector_probe.c
- kernel-prof.sh
- IntegrationSuite
- pathCache
- loadDaemonConfig
- filevault.go
- .netIntegrityContainer
- `[launch_rules]`
- LibraryClosure
- monitor
- IntegrationSuite
- absPath
- SyncSecretServiceFiles
- copyXattrs
- prctl_gate
- update/update.go
- .populateWatchInodes
- .awaitDaemonUpOrDead
- .TestRecord_GuardParity
- runNetworkGuard
- find_test.go
- linkView
- runLockdown
- inodeChain
- SanitizeText
- SampleContent
- supersede_probe.c
- launch_probe.c
- IntegrationSuite
- github.com/Virgula0/app-listener/tools/codecheck
- harnessSuite
- taint_victim.c
- secretJudge
- mutedText
- render.go
- Troubleshooting
- mountTable
- GenerateConf
- .runEditorHarness
- IntegrationSuite
- UnifiedDiff
- bpfLsmListed
- runGenKey
- auditGroupCoverage
- BinaryStat
- batchInodeRows
- configedit_test.go
- .eventTypeAllowed
- checkRunningBinaryMatchesInstalled
- guard_path_rename
- formatEntry
- bleedLayout
- Config
- parseKernelRelease
- golang.org/x/sys/unix.Timespec
- fileID
- MulticallError

## God Nodes (most connected - your core abstractions)
1. `Config` - 149 edges
2. `Guard` - 125 edges
3. `GuardInodeKey` - 113 edges
4. `Load()` - 101 edges
5. `absPath()` - 97 edges
6. `shell` - 97 edges
7. `IntegrationSuite` - 83 edges
8. `writeConfig()` - 79 edges
9. `fileEditModel` - 77 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Distributions` --references--> `monitor()`  [INFERRED]
  docs/compatibility.md → integrationtests/secret_probe/main.go
- `Files` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `Workspaces` --references--> `monitor()`  [INFERRED]
  docs/gui.md → integrationtests/secret_probe/main.go
- `eBPF at three levels` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go
- `Why not an AppArmor profile?` --references--> `monitor()`  [INFERRED]
  docs/how-it-works.md → integrationtests/secret_probe/main.go

## Import Cycles
- None detected.

## Communities (258 total, 45 thin omitted)

### Community 0 - "go_pkg_strings"
Cohesion: 0.11
Nodes (18): tempJournalRow, invalidPassword, recordTick, linkedFile, argvMatch, launchRule, trustEvent, btrfsIoctlFsInfoArgs (+10 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (182): plantBun(), TestOpenBunTmp_KeepsVettedDir(), newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestEveryInstallModeMapped() (+174 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (8): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), shQuote(), watchPathsInConfig()

### Community 3 - "SecretStoreNotes"
Cohesion: 0.39
Nodes (7): dockerKeyringHelper(), gitCredentialHelpers(), readOwnedRegular(), SecretStoreNotes(), TestReadOwnedRegular_RefusesSymlinkAndOtherOwner(), TestSecretStoreNotes(), WarnSecretStores()

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (85): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+77 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.23
Nodes (5): guardBinaryFlag(), IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent()

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (23): addErrKind, ReplacementCheck, classifyAddErr(), commMatchesGuardedBinary(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest (+15 more)

### Community 9 - "mkdirs"
Cohesion: 0.13
Nodes (37): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetHomeDepartureTriggersRefresh(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds() (+29 more)

### Community 10 - ".dispatch"
Cohesion: 0.09
Nodes (19): bpfGuardEvent, BusClient, AttestBusClient(), busClientKey(), execStamp(), imageOf(), pidfdPid(), procIdentity() (+11 more)

### Community 11 - "ConfinedEntry"
Cohesion: 0.14
Nodes (27): HashStore, sharedHashEntry, ConfinedEntry(), claimHash(), hashFileShared(), hashPathShared(), HashStats(), publishHash() (+19 more)

### Community 12 - "RefreshSection"
Cohesion: 0.05
Nodes (44): RefreshOptions, SectionScan, IsScript(), TestIsScript(), BinaryRule, homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob() (+36 more)

### Community 13 - "Guard"
Cohesion: 0.04
Nodes (38): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), binaryVerifyState, deferredBinary, exeRow, keptEntries (+30 more)

### Community 14 - "buildOneGuard"
Cohesion: 0.05
Nodes (59): DrainUntilSignal(), CheckBPFLSM(), CheckEBPF(), multicallAdmissible(), backingDeviceUnion(), buildGuards(), buildOneGuard(), eventFilterOptions() (+51 more)

### Community 15 - "time.Time"
Cohesion: 0.06
Nodes (44): LockedError, ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI() (+36 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (66): argv_env(), argv_matched(), argv_not_root(), argv_rule_at(), argv_short(), argv_word(), caller_updates(), creds_overridden() (+58 more)

### Community 18 - "navItem"
Cohesion: 0.04
Nodes (10): appTheme, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pillRenderer, tabButton, tabButtonRenderer (+2 more)

### Community 19 - "fileEditModel"
Cohesion: 0.07
Nodes (13): newFindInput(), clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser (+5 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.13
Nodes (52): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+44 more)

### Community 22 - "Resource"
Cohesion: 0.05
Nodes (14): fakeDaemon, Resource, narrowInLockedVault(), GuardRepository, Vault, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase (+6 more)

### Community 23 - "auth.go"
Cohesion: 0.10
Nodes (46): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+38 more)

### Community 24 - "op"
Cohesion: 0.17
Nodes (13): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+5 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.08
Nodes (79): inode_dev(), add_inode_to_guard(), admit_port(), btrfs_copy_gate(), bus_key(), check_and_emit_args(), chmod_only_drops_write(), code_suspect() (+71 more)

### Community 26 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.09
Nodes (30): class, findButton, pos, span, classOf(), column(), dedent(), Edit() (+22 more)

### Community 27 - "inode_dev.h"
Cohesion: 0.60
Nodes (3): ctx_ptr(), sb_dev(), guard_inode_free()

### Community 28 - "wm.go"
Cohesion: 0.11
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 29 - "rcUser"
Cohesion: 0.06
Nodes (70): addedResources(), TestAddedResourcesSkipsAlreadyGuardedSSH(), installFacultatives(), installSSHAgentFacultative(), listFacultatives(), runFacultatives(), facultativeUser(), TestInstallSSHAgentFacultative() (+62 more)

### Community 30 - "editView"
Cohesion: 0.16
Nodes (5): editView, Entry, fileSize(), joinRel(), relOf()

### Community 31 - "Network"
Cohesion: 0.09
Nodes (10): netGuardPill(), NetworkGuardRepository, NetworkMonitorRepository, clonePolicy(), NetEngines, NetExtra, NetPolicy, Network (+2 more)

### Community 32 - "CandidateDir"
Cohesion: 0.06
Nodes (48): appendSectionsAndEdit(), collectDiffAdditions(), TestDiffMergePreservesExistingAndParses(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories(), editConfig() (+40 more)

### Community 33 - "controlServer"
Cohesion: 0.16
Nodes (10): authSelfPeer(), readClientLines(), startControlServer(), streamEvents(), TestEditControlSessionEndIdempotent(), TestRecordFailureLockout(), configEditor, controlServer (+2 more)

### Community 34 - "picker"
Cohesion: 0.08
Nodes (31): keyEntry, listOpts, pickEntry, picker, pickMode, treeSearch, openPicker(), pickerTree() (+23 more)

### Community 35 - "secret_probe/main.go"
Cohesion: 0.18
Nodes (31): busPath(), connexec(), connFD(), createItem(), dial(), dupkey(), execInto(), failf() (+23 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (13): checkSecret(), checkTrustLine(), trustRows(), TestLiveSessionWatchAndPinger(), holdForward(), pingInterval(), backend, backendOp (+5 more)

### Community 38 - "systemd.go"
Cohesion: 0.07
Nodes (52): offerDaemonStart(), applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall(), installBinaryOnly(), markInstallComplete(), runRefreshCatalogOnly(), softenAutomatedRefreshErr() (+44 more)

### Community 39 - "NetEventType"
Cohesion: 0.08
Nodes (16): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), BinariesSummary(), TestBinariesSummary(), FormatAddr() (+8 more)

### Community 40 - "desktop.go"
Cohesion: 0.04
Nodes (44): releaseBuild(), Run(), codeEntry, installedBinary, latestRelease, pageDef, pageID, stageItem (+36 more)

### Community 41 - "UI"
Cohesion: 0.08
Nodes (20): nextConfig(), confirmUpdate(), isTerminal(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), cancelUI, EditText(), TestEditTextMapsUICancel() (+12 more)

### Community 42 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 44 - "pinstate.go"
Cohesion: 0.18
Nodes (22): ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID() (+14 more)

### Community 45 - "cursor.go"
Cohesion: 0.08
Nodes (20): buttonCursor, handCheck, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn() (+12 more)

### Community 46 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.04
Nodes (41): holdHeadless(), recordTickCmd(), HeadlessLine(), recordModel, snapGuard, diffModel, editorModel, GuardEvent (+33 more)

### Community 48 - "runView"
Cohesion: 0.05
Nodes (20): installUI, logGrid, logHook, logLine, multiAnswer, resizeWatch, runView, sizedDialog (+12 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (29): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+21 more)

### Community 50 - "EventType"
Cohesion: 0.07
Nodes (20): selfEventsFor(), snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting() (+12 more)

### Community 51 - "shell"
Cohesion: 0.04
Nodes (31): File, Progress, TrustItem, User, backendLink, configView, forwardForm, forwardView (+23 more)

### Community 52 - "editbackend_test.go"
Cohesion: 0.07
Nodes (20): Client, childEnv(), launch(), Start(), verifiedBinary(), helper(), TestChildEnvIsAllowlisted(), TestCleanRel() (+12 more)

### Community 53 - "changelog.go"
Cohesion: 0.11
Nodes (29): newChangelogModel(), showChangelog(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation() (+21 more)

### Community 54 - "Gate"
Cohesion: 0.16
Nodes (8): placedKey(), Caller, Gate, Role, Listen(), withSchema(), Bus, Judge

### Community 55 - ".apply"
Cohesion: 0.13
Nodes (11): Apply(), ProgressFunc, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal(), testPlan(), applier (+3 more)

### Community 56 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 57 - "serve.go"
Cohesion: 0.06
Nodes (36): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+28 more)

### Community 58 - "GuardInodeKey"
Cohesion: 0.03
Nodes (60): engine, portOwner, refusedReplacements, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), ensureBtrfsLayout() (+52 more)

### Community 59 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 60 - "FileTarget"
Cohesion: 0.14
Nodes (10): Target, Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps() (+2 more)

### Community 61 - "runMonitor"
Cohesion: 0.19
Nodes (14): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+6 more)

### Community 62 - "runDaemon"
Cohesion: 0.12
Nodes (16): newControlManager(), startControlManager(), makeReloadHandler(), runDaemon(), startPprof(), writePidFile(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted() (+8 more)

### Community 63 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.13
Nodes (20): placedWalk, treeJudge, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedEntry(), placedFileStat() (+12 more)

### Community 64 - "gui/session.go"
Cohesion: 0.08
Nodes (27): InstallMode, InstallRequest, Options, UninstallHooks, UninstallRequest, UpdateHooks, Run(), savedVariant() (+19 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 67 - "wire.go"
Cohesion: 0.16
Nodes (14): TestClientAuth(), TestSASL(), clientAuth(), decodeFrame(), distinctDictKeys(), header(), sameDict(), sameEntries() (+6 more)

### Community 68 - "buildTrustedSet"
Cohesion: 0.12
Nodes (20): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), allowedEnv(), binaryLaunchRules(), launchRulesByPath(), TestAllowedEnvKeepsOnlyRootPlaced(), TestLaunchRulesByPath(), buildTrustedSet() (+12 more)

### Community 69 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (6): IntegrationSuite, IntegrationSuite, IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite

### Community 70 - "Vault"
Cohesion: 0.11
Nodes (21): deprovisionKind, isFileVaultCiphertext(), isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+13 more)

### Community 71 - "expandPlaceholders"
Cohesion: 0.18
Nodes (7): symlinkTargetGlob(), expandPlaceholders(), TestExtraWatchPathsForSkipsMissing(), LibDirGlob, TrustGlob, reservedGlobs(), splitTrustGlob()

### Community 72 - "runNetworkMonitor"
Cohesion: 0.13
Nodes (11): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString() (+3 more)

### Community 73 - "fcntl"
Cohesion: 0.10
Nodes (9): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), denied() (+1 more)

### Community 74 - "policyEditor"
Cohesion: 0.13
Nodes (11): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+3 more)

### Community 75 - "gui.go"
Cohesion: 0.16
Nodes (17): runDiffCatalog(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary(), TestRunWithRefusesNonHeadlessBinary(), stagedBinary (+9 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "NetGuard"
Cohesion: 0.11
Nodes (10): foreignPinnedLink(), PumpRingbuf(), eventTypeKey(), judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), hook (+2 more)

### Community 78 - "liveSession"
Cohesion: 0.24
Nodes (9): parseErr(), TestLiveSessionConfigAndPut(), chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), normalizedEvents(), runForward() (+1 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 81 - "bufio.Reader"
Cohesion: 0.18
Nodes (13): readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readBareRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+5 more)

### Community 82 - "daemon.conf reference"
Cohesion: 0.12
Nodes (15): A complete example, `allow_lib`, Binary lines and events, daemon.conf reference, `[electron_apps]`, File format, `[inspectors]`, `lib_binary` (+7 more)

### Community 83 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 84 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (19): forwardCheck(), binaryStatOf(), TempBinary, TemporaryGrant, TempRule, ResolveTempBinary(), TestResolveTempBinaryClassesRuntime(), applyTempGrants() (+11 more)

### Community 85 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 86 - "exe_supersede.h"
Cohesion: 0.16
Nodes (19): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+11 more)

### Community 87 - "loginServer"
Cohesion: 0.09
Nodes (22): newCatalogRefresher(), markReadyWithLogin(), startCatalogRefresh(), hasSentinel(), lockParent(), makeRuntimeDir(), openUserRuntimeDir(), plantSentinel() (+14 more)

### Community 88 - "daemonPage"
Cohesion: 0.22
Nodes (3): daemonPage, daemonState, runHost

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LaunchRule"
Cohesion: 0.08
Nodes (34): ShortOpts, argvEntry, LaunchRule, LaunchWordForms(), EncodeLaunchRules(), ShortRefused(), canonLaunchRules(), encodeCatalog() (+26 more)

### Community 92 - "User"
Cohesion: 0.05
Nodes (57): fileID(), openBunTmp(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), entryWriters(), reserveChain(), TestGlobBuilder_LibBitsPerEntry(), underResource() (+49 more)

### Community 93 - "netinject.c"
Cohesion: 0.12
Nodes (15): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+7 more)

### Community 94 - "README.md"
Cohesion: 0.21
Nodes (8): Daemon configuration, Events, Flags, monitor: see what touches a path, Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 95 - "dialLiveSession"
Cohesion: 0.08
Nodes (27): auditAfterEdit(), auditFindings(), reserveStdout(), runBackend(), dialLiveSession(), LiveModeAvailable(), putConfig(), requireRoot() (+19 more)

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
Cohesion: 0.14
Nodes (14): procStats, statsSample, duration(), cpuTime(), humanBytes(), rssBytes(), New(), recv() (+6 more)

### Community 100 - "recordFlow"
Cohesion: 0.06
Nodes (41): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), formatElapsed(), lines(), offerBackupDeletion() (+33 more)

### Community 101 - "newTestTarget"
Cohesion: 0.09
Nodes (23): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+15 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "main_test.go"
Cohesion: 0.13
Nodes (6): IntegrationSuite, infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

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
Cohesion: 0.06
Nodes (48): entryOf(), readAtNoFollow(), upsertUnitAt(), writeAtNoFollow(), createOwned(), removeCheckedFile(), InstallCert, ReadInstallCert() (+40 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "planUpdaters"
Cohesion: 0.17
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 110 - "releasesSince"
Cohesion: 0.28
Nodes (11): newReleaseWindow(), releasesSince(), compareStableVersions(), latestTime(), newerThanStable(), parseStableVersion(), parseVerPart(), TestNewerThanStable() (+3 more)

### Community 111 - "run.go"
Cohesion: 0.36
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 112 - "SuperblockDev"
Cohesion: 0.40
Nodes (5): FirstOnBtrfs(), OnBtrfs(), SuperblockDev(), TestSuperblockDevMatchesStatOffBtrfs(), vouchedFd()

### Community 113 - "validateConfigText"
Cohesion: 0.11
Nodes (15): collectFilesystemPrereqs(), askEncryption(), discordCatalogEntry(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestSectionsFromCandidatesSkipsMissingWatchSubPaths() (+7 more)

### Community 115 - "catalogRefresher"
Cohesion: 0.09
Nodes (20): drain(), logRefreshChange(), readEvents(), swapFile(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails() (+12 more)

### Community 116 - "Row"
Cohesion: 0.29
Nodes (10): candidates(), confirmAll(), describe(), record(), run(), Given(), Row, Inspect() (+2 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "NetGuardEvent"
Cohesion: 0.20
Nodes (5): snapNetGuard, reasonText(), NetGuardEvent, TestReasonText(), netGuardEventLine

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "watchPattern"
Cohesion: 0.31
Nodes (9): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), inHome(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors() (+1 more)

### Community 124 - "inet_connect.c"
Cohesion: 0.25
Nodes (17): attack_connect(), attack_fastopen(), attack_uring(), denied(), dial(), fill_addr(), handoff(), listen_on() (+9 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "watchSet"
Cohesion: 0.19
Nodes (8): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, `[secret <name>]`, Sections at a glance

### Community 127 - "protected_test.go"
Cohesion: 0.19
Nodes (15): fakeLink, isOwnExecutable(), newFakeLink(), protectedShell(), TestForwardRefusesOwnExecutable(), TestProtectedEditFlow(), TestProtectedLiveCardsGated(), TestRecordReviewExtendsAnExistingSection() (+7 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "SelectOrphans"
Cohesion: 0.24
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 130 - "Recorder"
Cohesion: 0.14
Nodes (18): newBinary(), drain(), Recorder, Target, monitorTargets(), newRecorder(), opNames(), selfExe() (+10 more)

### Community 131 - "github.com/spf13/cobra.Command"
Cohesion: 0.07
Nodes (39): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+31 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "ValidateTargets"
Cohesion: 0.15
Nodes (13): existingSection(), guardedMatch(), overlaps(), refusedPath(), sectionFor(), validateTarget(), ValidateTargets(), mustParse() (+5 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "session"
Cohesion: 0.11
Nodes (7): connectionGone(), peerCred(), TestConnectionGoneOnlyFromTheBus(), readFrame(), queued, senderQueue, session

### Community 143 - "KernelDev"
Cohesion: 0.13
Nodes (19): fdStat, physicalParent(), statFD(), syncBusBrokers(), TagRelayBus(), KernelDev(), BtrfsDevices(), btrfsSysfsDevices() (+11 more)

### Community 144 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 145 - "Ledger"
Cohesion: 0.11
Nodes (5): applySchema(), Ledger, Pending, Source, fakeLedger

### Community 146 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.12
Nodes (16): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle() (+8 more)

### Community 147 - "maintainerKeeper"
Cohesion: 0.18
Nodes (12): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp (+4 more)

### Community 148 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 149 - "github.com/stretchr/testify/suite.Suite"
Cohesion: 0.20
Nodes (3): mountinfoUnitTest, unescapeMountinfo(), editorHarness

### Community 150 - "decryptStep"
Cohesion: 0.22
Nodes (9): askDecrypt(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty(), TestExistingEncryptedDirsNilPaths(), TestExistingEncryptedDirsSkipsPlain(), ExistingEncryptedDirs() (+1 more)

### Community 152 - "uninstallAll"
Cohesion: 0.12
Nodes (26): deletePostBackups(), listBackups(), restoreBackups(), confirmKeyDeletion(), offerBackupCleanup(), run(), runUninstall(), RunWith() (+18 more)

### Community 154 - "maintainer/plan.go"
Cohesion: 0.31
Nodes (7): isWhiteout(), opaque(), pathByte(), pathRecord(), TestPathRecord(), opKind, walker

### Community 155 - "guardLSMHooks"
Cohesion: 0.25
Nodes (5): lsmHook, trustHook, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), TestSocketHooksAreRequired()

### Community 156 - "io/fs.FileInfo"
Cohesion: 0.36
Nodes (9): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), fchownFromInfo(), futimesFromInfo(), readDirNames(), readlinkFD() (+1 more)

### Community 157 - "classify"
Cohesion: 0.14
Nodes (19): call(), dictFrame(), TestClassify(), TestClassifyCreateItemReplace(), TestClassifyItemTypeIsSchema(), TestClassifyOpenCallsExact(), TestDuplicateDictKeyRefused(), TestReadFrame() (+11 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "netGuardModel"
Cohesion: 0.09
Nodes (18): computeGuardBinaries(), Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero() (+10 more)

### Community 161 - "layout.go"
Cohesion: 0.32
Nodes (7): adaptiveGrid, checkGrid, gapVBox, cellMin(), newCheckGrid(), rowHeights(), visible()

### Community 162 - "update_test.go"
Cohesion: 0.06
Nodes (30): downloadFile(), newerThanInstalled(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestAssetsFor(), TestDownloadFileHTTPError() (+22 more)

### Community 163 - "SystemTrusted"
Cohesion: 0.29
Nodes (8): rootOwnedInode(), rootOwnedSafe(), SystemTrusted(), systemTool(), archiveMembers(), decompressed(), tarMembers(), TestArchiveMembersIgnoresNonArchives()

### Community 164 - "mc_tag.h"
Cohesion: 0.27
Nodes (9): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netm_exec_applet() (+1 more)

### Community 165 - "icon"
Cohesion: 0.04
Nodes (53): actionBar, binaryList, cell, column, counterStrip, dataTable, eventFeed, feedAction (+45 more)

### Community 166 - "TestSnapshots"
Cohesion: 0.08
Nodes (32): anyDaemonRunning(), reel, snapEngines, snapNet, sysStatus, newLogo(), newDaemonPage(), demoFiles() (+24 more)

### Community 167 - "secureResources"
Cohesion: 0.33
Nodes (7): encryptDirectories(), secureResources(), verifyEncryptionState(), confirmRunPrereq(), resolveFilesystemPrereqs(), EnsureSystemSetup(), isKernelAtLeast54()

### Community 168 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 169 - "sources_test.go"
Cohesion: 0.38
Nodes (4): newSources(), TestSourcesHoldFilesAndArchiveMembers(), TestSourcesRefuseOtherFilesystem(), writeTar()

### Community 170 - "`[watch <path>]`"
Cohesion: 0.33
Nodes (6): Installer placeholder, `[narrow <path>]`, `need_encryption`, `runtime:`, `watch:`, `[watch <path>]`

### Community 171 - "os.FileInfo"
Cohesion: 0.19
Nodes (5): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata()

### Community 172 - "Limitations"
Cohesion: 0.17
Nodes (12): Already-running processes, Binaries you can write, Commands that stop working, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root (+4 more)

### Community 173 - ".keyringObject"
Cohesion: 0.24
Nodes (3): placeableAttributes(), TestLearnKeyringFromFirstReply(), TestPlaceableAttributes()

### Community 174 - "plan_test.go"
Cohesion: 0.40
Nodes (11): scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles(), TestScanUntouchedLibDir(), TestTouches() (+3 more)

### Community 176 - "net_tester/main.go"
Cohesion: 0.19
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 177 - "Open"
Cohesion: 0.09
Nodes (31): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+23 more)

### Community 178 - "sources"
Cohesion: 0.19
Nodes (5): TestGeneratedLibraryBlockRoundTrips(), lstatOf(), outsideHardlinks(), libDirDev(), sources

### Community 179 - "tempGrantSetup"
Cohesion: 0.23
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 181 - "filevault_test.go"
Cohesion: 0.18
Nodes (26): looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault(), sealFileVault(), inode(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsEncryptedDispatchesToFileVaultForRegularFiles() (+18 more)

### Community 182 - "TrustGuard"
Cohesion: 0.05
Nodes (27): GlobChild, statFunc, tempGrant, tempMaskOp, tempRow, trustRows, deleteResKeys(), setSupersedeTrust() (+19 more)

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
Nodes (72): checkFreshDir(), openDirNoFollow(), renewDir(), setConfinedHomes(), VettedInode, Multicall, MulticallKind, bunAutoloadsBunfig() (+64 more)

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "sync.Mutex"
Cohesion: 0.20
Nodes (7): newActivity(), TestRingKeepsNewest(), NewNetwork(), Counts, newRing(), activity, ring

### Community 194 - "WithBottomBar"
Cohesion: 0.31
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

### Community 201 - "loadDaemonConfig"
Cohesion: 0.12
Nodes (17): loadDaemonConfig(), markReady(), noteReloadedSettings(), notifySystemdReady(), ensureAnchor(), ensurePortAnchors(), ensureRootDir(), ephemeralRange() (+9 more)

### Community 202 - "filevault.go"
Cohesion: 0.22
Nodes (14): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+6 more)

### Community 204 - "`[launch_rules]`"
Cohesion: 0.40
Nodes (5): `allow_env`, `deny_args`, `warn_args`, `deny_env`, `warn_env`, `[launch_rules]`, refuses()

### Community 205 - "LibraryClosure"
Cohesion: 0.09
Nodes (22): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths(), LibraryClosure() (+14 more)

### Community 206 - "monitor"
Cohesion: 0.06
Nodes (39): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks (+31 more)

### Community 209 - "absPath"
Cohesion: 0.05
Nodes (9): IntegrationSuite, IntegrationSuite, IntegrationSuite, IntegrationSuite, inNosuidNamespace(), IntegrationSuite, launchCase, absPath() (+1 more)

### Community 210 - "SyncSecretServiceFiles"
Cohesion: 0.50
Nodes (5): removeSessionHooks(), RemoveLoginDropIn(), installFile(), removeFile(), SyncSecretServiceFiles()

### Community 211 - "copyXattrs"
Cohesion: 0.50
Nodes (4): TestCopyXattrsPreservesUserAttributes(), copyXattrs(), listXattrNames(), splitXattrNames()

### Community 212 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 213 - "update/update.go"
Cohesion: 0.12
Nodes (43): installHooks(), staged(), changelogText(), cumulativeChangelog(), applyUpdate(), archAssets(), assetsFor(), channelOf() (+35 more)

### Community 216 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 217 - "runNetworkGuard"
Cohesion: 0.11
Nodes (21): confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), WithSuspectsFrom(), DiscoverInfraBinaries() (+13 more)

### Community 219 - "find_test.go"
Cohesion: 0.19
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 220 - "linkView"
Cohesion: 0.43
Nodes (3): linkTarget(), entryKind, linkView

### Community 221 - "runLockdown"
Cohesion: 0.14
Nodes (16): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), resolveConfigPath(), runLockdown(), dirIsRootOwnedSafe(), isBpffs() (+8 more)

### Community 222 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 223 - "SanitizeText"
Cohesion: 0.06
Nodes (21): newBinaryVetter(), openBinaryVetter(), ownerPaths(), scriptIdentity(), newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), binaryVetter (+13 more)

### Community 224 - "SampleContent"
Cohesion: 0.11
Nodes (24): TestDaemonUnitWithMetadataOutput(), installFileAs(), installServices(), installUnit(), upsertFile(), TestDaemonUnitCmdPrintsBundledUnit(), SampleContent(), SampleFiles() (+16 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 228 - "launch_probe.c"
Cohesion: 0.46
Nodes (6): anchor(), copy(), is_verb(), main(), read_once(), usage()

### Community 229 - "IntegrationSuite"
Cohesion: 0.25
Nodes (7): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests, IntegrationSuite

### Community 233 - "secretJudge"
Cohesion: 0.15
Nodes (8): secretJudge, SecretScope, SecretMatch, SecretService, parseSecretMatch(), TestSecretMatch(), LiveGuardAt(), Scope

### Community 234 - "mutedText"
Cohesion: 0.05
Nodes (41): actionGate, binaryRow, binSource, clickPanel, homePage, iconTile, installAction, installPage (+33 more)

### Community 235 - "render.go"
Cohesion: 0.13
Nodes (13): lineStyles, viewportState, baseStyle(), gutterLabel(), lineLen(), markStyle(), repeatSpaces(), setGutter() (+5 more)

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "mountTable"
Cohesion: 0.27
Nodes (7): MountInfo, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), userMountable(), mountTable

### Community 238 - "GenerateConf"
Cohesion: 0.05
Nodes (65): groupedDiscordConf(), TestGroupedSectionAddressedByEncryptionRoot(), TestPatchCatalogSectionGroupedConfig(), TestSetSectionWhitelistPreservesGroupStructure(), LibraryBlock, NarrowSection, confSafeMatch(), SecretServiceBlock() (+57 more)

### Community 241 - "UnifiedDiff"
Cohesion: 0.21
Nodes (12): applyEdit(), checkAndApply(), editAgain(), fetchConfig(), runEditConfig(), refreshCatalogConfig(), installConfig(), ConfirmOverwrite() (+4 more)

### Community 243 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 246 - "runGenKey"
Cohesion: 0.20
Nodes (11): confirmMasterKeyOverwrite(), runGenKey(), runOneShotMode(), ensureInstalledBinary(), ensureMasterKey(), prepareInstallation(), TestEnsureInstalledBinary(), checkKeyLen() (+3 more)

### Community 247 - "auditGroupCoverage"
Cohesion: 0.22
Nodes (11): auditEntry(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths(), relOrPath() (+3 more)

### Community 249 - "BinaryStat"
Cohesion: 0.21
Nodes (13): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+5 more)

### Community 252 - "batchInodeRows"
Cohesion: 0.27
Nodes (10): inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows(), walkInodeBatches(), countInodeRows() (+2 more)

### Community 254 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 257 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 258 - "checkRunningBinaryMatchesInstalled"
Cohesion: 0.40
Nodes (5): checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), fileSHA256(), sameContent(), TestCheckRunningBinaryMatchesInstalled()

### Community 260 - "guard_path_rename"
Cohesion: 0.27
Nodes (18): get_dentry_from_path(), get_inode_from_path(), guard_inode_getattr(), guard_path_link(), guard_path_mkdir(), guard_path_mknod(), guard_path_rename(), guard_path_rmdir() (+10 more)

### Community 262 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 267 - "Config"
Cohesion: 0.06
Nodes (79): sameConfig(), watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyBinaryLine(), applyBlockLine(), applyConfigLine() (+71 more)

### Community 271 - "parseKernelRelease"
Cohesion: 0.67
Nodes (4): ciliumVersion, parseKernelRelease(), TestParseKernelRelease(), unameKernelVersion()

## Knowledge Gaps
- **34 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck` (+29 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 691 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **45 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `monitor()` connect `monitor` to `time.Duration`, `secret_probe/main.go`, `Installation`?**
  _High betweenness centrality (0.031) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _34 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `go_pkg_strings` be split into smaller, more focused modules?**
  _Cohesion score 0.10734177215189873 - nodes in this community are weakly interconnected._
- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `parseDaemonEvents`, `main_test.go`, `IntegrationSuite`, `.netIntegrityContainer`, `IntegrationSuite`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `.TestRecord_GuardParity`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.01602979274611399 - nodes in this community are weakly interconnected._
- **Why does `engine` connect `GuardInodeKey` to `go_pkg_strings`, `sync.Mutex`, `NetGuard`, `buildOneGuard`, `Guard`, `guardLSMHooks`, `batchInodeRows`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.0943859649122807 - nodes in this community are weakly interconnected._